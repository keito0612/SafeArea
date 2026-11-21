package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/isobekeito/safearea/internal/domain/model"
	"github.com/isobekeito/safearea/internal/mock"
)

// テスト用のセットアップヘルパー
func setupDisasterHistoryUseCase(t *testing.T) (*gomock.Controller, *mock.MockDisasterRepository, DisasterHistoryUseCase) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockRepo := mock.NewMockDisasterRepository(ctrl)
	uc := NewDisasterHistoryUseCase(mockRepo)
	return ctrl, mockRepo, uc
}

func TestDisasterHistoryUseCase_GetDisasterHistory(t *testing.T) {
	t.Parallel()

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name          string
			address       string
			mockSetup     func(*mock.MockDisasterRepository)
			wantCount     int
			wantFirstID   int64
			wantFirstType string
		}{
			{
				name:    "市区町村で検索成功",
				address: "東京都渋谷区神南1丁目",
				mockSetup: func(m *mock.MockDisasterRepository) {
					m.EXPECT().
						FindByCity(gomock.Any(), "13", "渋谷区").
						Return([]*model.Disaster{
							{
								ID:     1,
								Type:   model.DisasterTypeEarthquake,
								Name:   "東日本大震災",
								Year:   2011,
								Damage: "液状化による被害",
							},
						}, nil)
				},
				wantCount:     1,
				wantFirstID:   1,
				wantFirstType: "earthquake",
			},
			{
				name:    "都道府県のみで検索成功",
				address: "東京都",
				mockSetup: func(m *mock.MockDisasterRepository) {
					m.EXPECT().
						FindByPrefecture(gomock.Any(), "13").
						Return([]*model.Disaster{
							{ID: 1, Type: model.DisasterTypeFlood, Name: "台風19号", Year: 2019},
						}, nil)
				},
				wantCount:     1,
				wantFirstID:   1,
				wantFirstType: "flood",
			},
			{
				name:    "災害履歴が空",
				address: "東京都渋谷区",
				mockSetup: func(m *mock.MockDisasterRepository) {
					m.EXPECT().
						FindByCity(gomock.Any(), "13", "渋谷区").
						Return([]*model.Disaster{}, nil)
				},
				wantCount: 0,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				// Arrange
				ctrl, mockRepo, uc := setupDisasterHistoryUseCase(t)
				defer ctrl.Finish()
				tt.mockSetup(mockRepo)

				// Act
				result, err := uc.GetDisasterHistory(context.Background(), tt.address)

				// Assert
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Len(t, result.Disasters, tt.wantCount)

				if tt.wantCount > 0 {
					assert.Equal(t, tt.wantFirstID, result.Disasters[0].ID)
					assert.Equal(t, tt.wantFirstType, result.Disasters[0].Type)
				}
			})
		}
	})

	t.Run("エラー系", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name           string
			address        string
			mockSetup      func(*mock.MockDisasterRepository)
			wantErrContain string
		}{
			{
				name:           "無効な住所",
				address:        "無効な住所",
				mockSetup:      func(m *mock.MockDisasterRepository) {},
				wantErrContain: "prefecture not found",
			},
			{
				name:    "リポジトリエラー",
				address: "東京都渋谷区",
				mockSetup: func(m *mock.MockDisasterRepository) {
					m.EXPECT().
						FindByCity(gomock.Any(), "13", "渋谷区").
						Return(nil, errors.New("database error"))
				},
				wantErrContain: "failed to find disasters",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				// Arrange
				ctrl, mockRepo, uc := setupDisasterHistoryUseCase(t)
				defer ctrl.Finish()
				tt.mockSetup(mockRepo)

				// Act
				result, err := uc.GetDisasterHistory(context.Background(), tt.address)

				// Assert
				require.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), tt.wantErrContain)
			})
		}
	})

	t.Run("データ変換の確認", func(t *testing.T) {
		t.Parallel()

		// Arrange
		ctrl, mockRepo, uc := setupDisasterHistoryUseCase(t)
		defer ctrl.Finish()

		mockRepo.EXPECT().
			FindByCity(gomock.Any(), "13", "渋谷区").
			Return([]*model.Disaster{
				{
					ID:     123,
					Type:   model.DisasterTypeEarthquake,
					Name:   "テスト地震",
					Year:   2020,
					Damage: "テスト被害",
				},
			}, nil)

		// Act
		result, err := uc.GetDisasterHistory(context.Background(), "東京都渋谷区")

		// Assert: モデルからDTOへの変換が正しいか
		require.NoError(t, err)
		require.Len(t, result.Disasters, 1)

		disaster := result.Disasters[0]
		assert.Equal(t, int64(123), disaster.ID)
		assert.Equal(t, "earthquake", disaster.Type) // model.DisasterTypeEarthquake → "earthquake"
		assert.Equal(t, "テスト地震", disaster.Name)
		assert.Equal(t, 2020, disaster.Year)
		assert.Equal(t, "テスト被害", disaster.Damage)
	})
}
