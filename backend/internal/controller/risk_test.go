package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/isobekeito/safearea/internal/domain/model"
	"github.com/isobekeito/safearea/internal/dto/response"
	"github.com/isobekeito/safearea/internal/infrastructure/external"
	"github.com/isobekeito/safearea/internal/mock"
	"github.com/isobekeito/safearea/internal/usecase"
)

// テスト用のセットアップヘルパー
func setupSearchController(t *testing.T) (*gomock.Controller, *mock.MockGeocodingClient, *mock.MockDisasterRepository, *SearchController) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockGeoClient := mock.NewMockGeocodingClient(ctrl)
	mockRepo := mock.NewMockDisasterRepository(ctrl)
	uc := usecase.NewSearchUseCase(mockGeoClient, mockRepo)
	controller := NewSearchController(uc)
	return ctrl, mockGeoClient, mockRepo, controller
}

// HTTPリクエスト実行ヘルパー
func executeSearchRequest(t *testing.T, controller *SearchController, url string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()
	controller.Search(rec, req)
	return rec
}

// ジオコーディング成功のモックセットアップ
func setupGeocodeMock(m *mock.MockGeocodingClient, address string, lat, lng float64, code string) {
	m.EXPECT().
		Geocode(gomock.Any(), address).
		Return(&external.GeocodingResult{
			Coordinates: external.Coordinates{Lat: lat, Lng: lng},
			AddressCode: code,
		}, nil)
}

func TestSearchController_Search(t *testing.T) {
	t.Parallel()

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name              string
			url               string
			address           string
			mockGeoSetup      func(*mock.MockGeocodingClient)
			mockRepoSetup     func(*mock.MockDisasterRepository)
			wantStatusCode    int
			wantDisasterCount float64
		}{
			{
				name:    "市区町村で検索",
				url:     "/api/search?address=東京都渋谷区",
				address: "東京都渋谷区",
				mockGeoSetup: func(m *mock.MockGeocodingClient) {
					setupGeocodeMock(m, "東京都渋谷区", 35.658, 139.701, "13113")
				},
				mockRepoSetup: func(m *mock.MockDisasterRepository) {
					m.EXPECT().
						FindByCity(gomock.Any(), "13", "渋谷区").
						Return([]*model.Disaster{
							{ID: 1, Name: "災害1"},
							{ID: 2, Name: "災害2"},
						}, nil)
				},
				wantStatusCode:    http.StatusOK,
				wantDisasterCount: 2,
			},
			{
				name:    "都道府県のみで検索",
				url:     "/api/search?address=東京都",
				address: "東京都",
				mockGeoSetup: func(m *mock.MockGeocodingClient) {
					setupGeocodeMock(m, "東京都", 35.689, 139.692, "13000")
				},
				mockRepoSetup: func(m *mock.MockDisasterRepository) {
					m.EXPECT().
						FindByPrefecture(gomock.Any(), "13").
						Return([]*model.Disaster{
							{ID: 1, Name: "災害1"},
							{ID: 2, Name: "災害2"},
							{ID: 3, Name: "災害3"},
						}, nil)
				},
				wantStatusCode:    http.StatusOK,
				wantDisasterCount: 3,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				// Arrange
				ctrl, mockGeoClient, mockRepo, controller := setupSearchController(t)
				defer ctrl.Finish()
				tt.mockGeoSetup(mockGeoClient)
				tt.mockRepoSetup(mockRepo)

				// Act
				rec := executeSearchRequest(t, controller, tt.url)

				// Assert
				assert.Equal(t, tt.wantStatusCode, rec.Code)

				var resp map[string]any
				err := json.NewDecoder(rec.Body).Decode(&resp)
				require.NoError(t, err)
				assert.True(t, resp["success"].(bool))

				data := resp["data"].(map[string]any)
				assert.Equal(t, tt.wantDisasterCount, data["disaster_history_count"])
			})
		}
	})

	t.Run("バリデーションエラー", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			url  string
		}{
			{name: "addressパラメータなし", url: "/api/search"},
			{name: "空のaddress", url: "/api/search?address="},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				// Arrange
				ctrl, _, _, controller := setupSearchController(t)
				defer ctrl.Finish()

				// Act
				rec := executeSearchRequest(t, controller, tt.url)

				// Assert
				assert.Equal(t, http.StatusBadRequest, rec.Code)

				var resp response.ErrorResponse
				err := json.NewDecoder(rec.Body).Decode(&resp)
				require.NoError(t, err)
				assert.False(t, resp.Success)
				assert.Equal(t, response.ErrorCodeValidation, resp.Error.Code)
			})
		}
	})

	t.Run("エラー系", func(t *testing.T) {
		t.Parallel()

		t.Run("ジオコーディングエラー", func(t *testing.T) {
			t.Parallel()

			// Arrange
			ctrl, mockGeoClient, _, controller := setupSearchController(t)
			defer ctrl.Finish()

			mockGeoClient.EXPECT().
				Geocode(gomock.Any(), "存在しない住所").
				Return(nil, errors.New("no results found for address"))

			// Act
			rec := executeSearchRequest(t, controller, "/api/search?address=存在しない住所")

			// Assert
			assert.Equal(t, http.StatusNotFound, rec.Code)

			var resp response.ErrorResponse
			err := json.NewDecoder(rec.Body).Decode(&resp)
			require.NoError(t, err)
			assert.Equal(t, response.ErrorCodeAddressNotFound, resp.Error.Code)
		})

		t.Run("リポジトリエラー時も正常レスポンス", func(t *testing.T) {
			t.Parallel()

			// Arrange
			ctrl, mockGeoClient, mockRepo, controller := setupSearchController(t)
			defer ctrl.Finish()

			setupGeocodeMock(mockGeoClient, "東京都渋谷区", 35.658, 139.701, "13113")
			mockRepo.EXPECT().
				FindByCity(gomock.Any(), "13", "渋谷区").
				Return(nil, errors.New("database error"))

			// Act
			rec := executeSearchRequest(t, controller, "/api/search?address=東京都渋谷区")

			// Assert: 災害履歴取得エラーでも座標は返せるので200 OK
			assert.Equal(t, http.StatusOK, rec.Code)

			var resp map[string]any
			err := json.NewDecoder(rec.Body).Decode(&resp)
			require.NoError(t, err)

			data := resp["data"].(map[string]any)
			assert.Equal(t, float64(0), data["disaster_history_count"])
		})
	})

	t.Run("レスポンス詳細確認", func(t *testing.T) {
		t.Parallel()

		// Arrange
		ctrl, mockGeoClient, mockRepo, controller := setupSearchController(t)
		defer ctrl.Finish()

		setupGeocodeMock(mockGeoClient, "東京都江東区", 35.655, 139.793, "13108")
		mockRepo.EXPECT().
			FindByCity(gomock.Any(), "13", "江東区").
			Return([]*model.Disaster{
				{ID: 1, Name: "台風19号"},
				{ID: 2, Name: "東日本大震災"},
			}, nil)

		// Act
		rec := executeSearchRequest(t, controller, "/api/search?address=東京都江東区")

		// Assert
		require.Equal(t, http.StatusOK, rec.Code)

		var resp map[string]any
		err := json.NewDecoder(rec.Body).Decode(&resp)
		require.NoError(t, err)

		data := resp["data"].(map[string]any)
		assert.Equal(t, "東京都江東区", data["address"])
		assert.Equal(t, float64(2), data["disaster_history_count"])

		coords := data["coordinates"].(map[string]any)
		assert.Equal(t, 35.655, coords["lat"])
		assert.Equal(t, 139.793, coords["lng"])
	})
}
