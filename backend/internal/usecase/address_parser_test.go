package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAddress(t *testing.T) {
	tests := []struct {
		name     string
		address  string
		expected *AddressInfo
	}{
		{
			name:    "東京都渋谷区",
			address: "東京都渋谷区神南1丁目",
			expected: &AddressInfo{
				Prefecture:     "東京都",
				PrefectureCode: "13",
				City:           "渋谷区",
			},
		},
		{
			name:    "北海道札幌市",
			address: "北海道札幌市中央区",
			expected: &AddressInfo{
				Prefecture:     "北海道",
				PrefectureCode: "01",
				City:           "札幌市",
			},
		},
		{
			name:    "大阪府大阪市",
			address: "大阪府大阪市北区梅田",
			expected: &AddressInfo{
				Prefecture:     "大阪府",
				PrefectureCode: "27",
				City:           "大阪市",
			},
		},
		{
			name:    "神奈川県横浜市",
			address: "神奈川県横浜市西区みなとみらい",
			expected: &AddressInfo{
				Prefecture:     "神奈川県",
				PrefectureCode: "14",
				City:           "横浜市",
			},
		},
		{
			name:    "都道府県のみ",
			address: "東京都",
			expected: &AddressInfo{
				Prefecture:     "東京都",
				PrefectureCode: "13",
				City:           "",
			},
		},
		{
			name:    "無効な住所",
			address: "無効な住所",
			expected: &AddressInfo{
				Prefecture:     "",
				PrefectureCode: "",
				City:           "",
			},
		},
		{
			name:    "村を含む住所",
			address: "長野県北安曇郡白馬村",
			expected: &AddressInfo{
				Prefecture:     "長野県",
				PrefectureCode: "20",
				City:           "北安曇郡白馬村",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAddress(tt.address)
			assert.Equal(t, tt.expected.Prefecture, result.Prefecture)
			assert.Equal(t, tt.expected.PrefectureCode, result.PrefectureCode)
			assert.Equal(t, tt.expected.City, result.City)
		})
	}
}
