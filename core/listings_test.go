package core_test

import (
	"testing"

	"github.com/OpenBazaar/openbazaar-go/core"
	"github.com/OpenBazaar/openbazaar-go/pb"
	"github.com/OpenBazaar/openbazaar-go/test/factory"
)

func TestFactoryCryptoListingCoinDivisibilityMatchesConst(t *testing.T) {
	if factory.NewCryptoListing("blu").Metadata.CoinDivisibility != core.DefaultCoinDivisibility {
		t.Fatal("DefaultCoinDivisibility constant has changed. Please update factory value.")
	}
}

func TestVersionForNewListing(t *testing.T) {
	listing := &pb.Listing{Metadata: &pb.Listing_Metadata{
		PriceModifier: 1,
		Format:        pb.Listing_Metadata_MARKET_PRICE,
	}}

	for i, test := range []struct {
		modifier        float32
		currentVersion  uint32
		format          pb.Listing_Metadata_Format
		expectedVersion uint32
	}{
		// Old version - no markup
		{0, core.PriceModifierListingVersion, pb.Listing_Metadata_MARKET_PRICE, core.PriceModifierListingVersion - 1},

		// Old version - wrong format
		{1, core.PriceModifierListingVersion, pb.Listing_Metadata_FIXED_PRICE, core.PriceModifierListingVersion - 1},

		// Old version - wrong current version
		{1, core.PriceModifierListingVersion - 1, pb.Listing_Metadata_MARKET_PRICE, core.PriceModifierListingVersion - 1},
		{1, core.PriceModifierListingVersion + 1, pb.Listing_Metadata_MARKET_PRICE, core.PriceModifierListingVersion + 1},
		{1, core.PriceModifierListingVersion + 7, pb.Listing_Metadata_MARKET_PRICE, core.PriceModifierListingVersion + 7},

		// New version
		{1, core.PriceModifierListingVersion, pb.Listing_Metadata_MARKET_PRICE, core.PriceModifierListingVersion},
		{-1, core.PriceModifierListingVersion, pb.Listing_Metadata_MARKET_PRICE, core.PriceModifierListingVersion},
	} {
		listing.Metadata.Format = test.format
		listing.Metadata.PriceModifier = test.modifier
		v := core.VersionForNewListingAndCurrentVersion(listing, test.currentVersion)
		if v != test.expectedVersion {
			t.Fatal("Test", i, "failed\nWanted:", test.expectedVersion, "\nGot:", v)
		}
	}
}
