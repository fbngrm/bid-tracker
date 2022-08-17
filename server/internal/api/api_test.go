package api_test

// import (
// 	"context"
// 	"errors"
// 	"testing"

// 	apiv1 "github.com/fbngrm/bid-tracker/gen/proto/go/match/v1"
// 	"github.com/fbngrm/bid-tracker/pkg/location"
// 	"github.com/fbngrm/bid-tracker/pkg/match"
// 	"github.com/fbngrm/bid-tracker/pkg/materials"
// 	"github.com/fbngrm/bid-tracker/pkg/partner"
// 	"github.com/fbngrm/bid-tracker/server/api"
// 	"github.com/fbngrm/bid-tracker/server/api/mocks"
// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/assert"
// )

// func TestMatchPartnersWithRequest(t *testing.T) {
// 	testcases := []struct {
// 		name    string
// 		req     *apiv1.MatchPartnersWithRequestInput
// 		resp    *apiv1.MatchPartnersWithRequestOutput
// 		err     error
// 		matches match.Matches
// 	}{
// 		{
// 			name: "expect 0 partner to be returned, no error",
// 			req: &apiv1.MatchPartnersWithRequestInput{
// 				Material: "wood",
// 				Location: "52.532566:13.396261",
// 			},
// 			resp: &apiv1.MatchPartnersWithRequestOutput{
// 				Partner: []*apiv1.Partner{},
// 			},
// 			err:     nil,
// 			matches: match.Matches{},
// 		},
// 		{
// 			name: "expect 1 partner to be returned, no error",
// 			req: &apiv1.MatchPartnersWithRequestInput{
// 				Material: "wood",
// 				Location: "52.532566:13.396261",
// 			},
// 			resp: &apiv1.MatchPartnersWithRequestOutput{
// 				Partner: []*apiv1.Partner{
// 					{
// 						Id:        uint32(1),
// 						Materials: "WOOD",
// 						Rating:    uint32(5),
// 						Location:  "52.500000:13.400000",
// 						Radius:    50000.,
// 						Distance:  10000.,
// 					},
// 				},
// 			},
// 			err: nil,
// 			matches: match.Matches{
// 				{
// 					Partner: partner.Partner{
// 						ID:        1,
// 						Materials: []materials.Material{"WOOD"},
// 						Address: location.Location{
// 							Lat:  52.500000,
// 							Long: 13.400000,
// 						},
// 						RadiusOfOperation: 50000.,
// 						Rating:            5,
// 					},
// 					Distance: 10000.,
// 				},
// 			},
// 		},
// 		{
// 			name: "expect 2 partner in same order, no error",
// 			req: &apiv1.MatchPartnersWithRequestInput{
// 				Material: "wood",
// 				Location: "52.532566:13.396261",
// 			},
// 			resp: &apiv1.MatchPartnersWithRequestOutput{
// 				Partner: []*apiv1.Partner{
// 					{
// 						Id:        uint32(2),
// 						Materials: "WOOD, CARPET",
// 						Rating:    uint32(4),
// 						Location:  "52.400000:13.300000",
// 						Radius:    50000.,
// 						Distance:  10000.,
// 					},
// 					{
// 						Id:        uint32(1),
// 						Materials: "WOOD",
// 						Rating:    uint32(5),
// 						Location:  "52.500000:13.400000",
// 						Radius:    50000.,
// 						Distance:  10000.,
// 					},
// 				},
// 			},
// 			err: nil,
// 			matches: match.Matches{
// 				{
// 					Partner: partner.Partner{
// 						ID:        2,
// 						Materials: []materials.Material{"WOOD", "CARPET"},
// 						Address: location.Location{
// 							Lat:  52.400000,
// 							Long: 13.300000,
// 						},
// 						RadiusOfOperation: 50000.,
// 						Rating:            4,
// 					},
// 					Distance: 10000.,
// 				},
// 				{
// 					Partner: partner.Partner{
// 						ID:        1,
// 						Materials: []materials.Material{"WOOD"},
// 						Address: location.Location{
// 							Lat:  52.500000,
// 							Long: 13.400000,
// 						},
// 						RadiusOfOperation: 50000.,
// 						Rating:            5,
// 					},
// 					Distance: 10000.,
// 				},
// 			},
// 		},
// 		{
// 			name: "expect error for non-existing material",
// 			req: &apiv1.MatchPartnersWithRequestInput{
// 				Material: "glass",
// 				Location: "52.532566:13.396261",
// 			},
// 			resp: &apiv1.MatchPartnersWithRequestOutput{},
// 			err:  errors.New("could not parse material: material not supported [\"glass\"]"),
// 		},
// 	}

// 	controller := gomock.NewController(t)
// 	matchService := mocks.NewMockMatcher(controller)
// 	api := api.NewApi(matchService)

// 	for _, tc := range testcases {
// 		if tc.err == nil {
// 			matchService.EXPECT().GetMatches(gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.matches, tc.err)
// 		}
// 		resp, err := api.MatchPartnersWithRequest(context.Background(), tc.req)

// 		// expected error
// 		if tc.err != nil {
// 			if err == nil {
// 				t.Log(tc.name)
// 				t.Logf("expected error but got nil\n")
// 				t.FailNow()
// 			}
// 			if !assert.EqualError(t, err, tc.err.Error()) {
// 				t.Log(tc.name, " failed")
// 			}
// 		}

// 		// unexpected error
// 		if err != nil && tc.err == nil {
// 			t.Log(tc.name)
// 			t.Logf("expected no error but got: %v\n", err)
// 			t.FailNow()
// 		}

// 		got := resp.GetPartner()
// 		want := tc.resp.GetPartner()
// 		if !assert.Equal(t, want, got) {
// 			t.Log(tc.name, " failed")
// 		}
// 	}
// }
