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
	"github.com/isobekeito/safearea/internal/mock"
	"github.com/isobekeito/safearea/internal/usecase"
)

// テスト用のセットアップヘルパー
func setupDisastersController(t *testing.T) (*gomock.Controller, *mock.MockDisasterRepository, *DisastersController) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockRepo := mock.NewMockDisasterRepository(ctrl)
	uc := usecase.NewDisasterHistoryUseCase(mockRepo)
	controller := NewDisastersController(uc)
	return ctrl, mockRepo, controller
}

// HTTPリクエスト実行ヘルパー
func executeDisastersRequest(t *testing.T, controller *DisastersController, url string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()
	controller.GetDisasters(rec, req)
	return rec
}

func TestDisastersController_GetDisasters(t *testing.T) {
	t.Parallel()

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name           string
			url            string
			mockSetup      func(*mock.MockDisasterRepository)
			wantStatusCode int
			wantCount      int
		}{
			{
				name: "複数の災害履歴を取得",
				url:  "/api/search/disasters?address=東京都渋谷区",
				mockSetup: func(m *mock.MockDisasterRepository) {
					m.EXPECT().
						FindByCity(gomock.Any(), "13", "渋谷区").
						Return([]*model.Disaster{
							{ID: 1, Type: model.DisasterTypeEarthquake, Name: "東日本大震災", Year: 2011},
							{ID: 2, Type: model.DisasterTypeFlood, Name: "台風19号", Year: 2019},
						}, nil)
				},
				wantStatusCode: http.StatusOK,
				wantCount:      2,
			},
			{
				name: "都道府県のみで検索",
				url:  "/api/search/disasters?address=東京都",
				mockSetup: func(m *mock.MockDisasterRepository) {
					m.EXPECT().
						FindByPrefecture(gomock.Any(), "13").
						Return([]*model.Disaster{
							{ID: 1, Type: model.DisasterTypeFlood, Name: "台風19号", Year: 2019},
						}, nil)
				},
				wantStatusCode: http.StatusOK,
				wantCount:      1,
			},
			{
				name: "災害履歴が空",
				url:  "/api/search/disasters?address=北海道札幌市",
				mockSetup: func(m *mock.MockDisasterRepository) {
					m.EXPECT().
						FindByCity(gomock.Any(), "01", "札幌市").
						Return([]*model.Disaster{}, nil)
				},
				wantStatusCode: http.StatusOK,
				wantCount:      0,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				// Arrange
				ctrl, mockRepo, controller := setupDisastersController(t)
				defer ctrl.Finish()
				tt.mockSetup(mockRepo)

				// Act
				rec := executeDisastersRequest(t, controller, tt.url)

				// Assert
				assert.Equal(t, tt.wantStatusCode, rec.Code)

				var resp map[string]any
				err := json.NewDecoder(rec.Body).Decode(&resp)
				require.NoError(t, err)
				assert.True(t, resp["success"].(bool))

				data := resp["data"].(map[string]any)
				disasters := data["disasters"].([]any)
				assert.Len(t, disasters, tt.wantCount)
			})
		}
	})

	t.Run("バリデーションエラー", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			url  string
		}{
			{name: "addressパラメータなし", url: "/api/search/disasters"},
			{name: "空のaddress", url: "/api/search/disasters?address="},
			{name: "空白のみのaddress", url: "/api/search/disasters?address=%20%20%20"},
			{name: "無効な都道府県", url: "/api/search/disasters?address=無効な住所"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				// Arrange
				ctrl, _, controller := setupDisastersController(t)
				defer ctrl.Finish()
				// Repositoryは呼ばれない（バリデーションで弾かれる）

				// Act
				rec := executeDisastersRequest(t, controller, tt.url)

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

		tests := []struct {
			name           string
			url            string
			mockSetup      func(*mock.MockDisasterRepository)
			wantStatusCode int
			wantErrorCode  response.ErrorCode
		}{
			{
				name: "リポジトリエラー",
				url:  "/api/search/disasters?address=東京都渋谷区",
				mockSetup: func(m *mock.MockDisasterRepository) {
					m.EXPECT().
						FindByCity(gomock.Any(), "13", "渋谷区").
						Return(nil, errors.New("database error"))
				},
				wantStatusCode: http.StatusInternalServerError,
				wantErrorCode:  response.ErrorCodeInternal,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				// Arrange
				ctrl, mockRepo, controller := setupDisastersController(t)
				defer ctrl.Finish()
				tt.mockSetup(mockRepo)

				// Act
				rec := executeDisastersRequest(t, controller, tt.url)

				// Assert
				assert.Equal(t, tt.wantStatusCode, rec.Code)

				var resp response.ErrorResponse
				err := json.NewDecoder(rec.Body).Decode(&resp)
				require.NoError(t, err)
				assert.False(t, resp.Success)
				assert.Equal(t, tt.wantErrorCode, resp.Error.Code)
			})
		}
	})

	t.Run("レスポンス詳細確認", func(t *testing.T) {
		t.Parallel()

		// Arrange
		ctrl, mockRepo, controller := setupDisastersController(t)
		defer ctrl.Finish()

		mockRepo.EXPECT().
			FindByCity(gomock.Any(), "13", "渋谷区").
			Return([]*model.Disaster{
				{
					ID:     1,
					Type:   model.DisasterTypeEarthquake,
					Name:   "東日本大震災",
					Year:   2011,
					Damage: "液状化による地盤沈下",
				},
			}, nil)

		// Act
		rec := executeDisastersRequest(t, controller, "/api/search/disasters?address=東京都渋谷区")

		// Assert
		require.Equal(t, http.StatusOK, rec.Code)

		var resp map[string]any
		err := json.NewDecoder(rec.Body).Decode(&resp)
		require.NoError(t, err)

		data := resp["data"].(map[string]any)
		disasters := data["disasters"].([]any)
		require.Len(t, disasters, 1)

		firstDisaster := disasters[0].(map[string]any)
		assert.Equal(t, float64(1), firstDisaster["id"])
		assert.Equal(t, "earthquake", firstDisaster["type"])
		assert.Equal(t, "東日本大震災", firstDisaster["name"])
		assert.Equal(t, float64(2011), firstDisaster["year"])
		assert.Equal(t, "液状化による地盤沈下", firstDisaster["damage"])
	})
}
