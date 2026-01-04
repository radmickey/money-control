package converters

import (
	"strings"

	assetspb "github.com/radmickey/money-control/backend/proto/assets"
)

// Asset type conversions for assets service
var assetTypeAssetsMap = map[string]assetspb.AssetType{
	"stock":       assetspb.AssetType_ASSET_TYPE_STOCK,
	"crypto":      assetspb.AssetType_ASSET_TYPE_CRYPTO,
	"etf":         assetspb.AssetType_ASSET_TYPE_ETF,
	"real_estate": assetspb.AssetType_ASSET_TYPE_REAL_ESTATE,
	"cash":        assetspb.AssetType_ASSET_TYPE_CASH,
	"bond":        assetspb.AssetType_ASSET_TYPE_BOND,
	"other":       assetspb.AssetType_ASSET_TYPE_OTHER,
}

func StringToAssetTypeAssets(s string) assetspb.AssetType {
	if t, ok := assetTypeAssetsMap[strings.ToLower(s)]; ok {
		return t
	}
	return assetspb.AssetType_ASSET_TYPE_OTHER
}

