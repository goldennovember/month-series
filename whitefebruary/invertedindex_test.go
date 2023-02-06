package whitefebruary

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var SmallFile = "./resources/small_abstracts.json"
var BigFile = "./resources/big_abstracts.json"

func Test_smallGetPageUrlsForTerm(t *testing.T) {
	ast := assert.New(t)
	smallFileInvertedIndex := InvertedIndex{Filename: SmallFile, HashMap: map[string]*Data{}}
	smallFileInvertedIndex.buildIndex()

	t.Run("smallGetSingleTermOrganisms", func(t *testing.T) {
		got := smallFileInvertedIndex.getPageUrlsForTerm([]string{"organisms"})

		fmt.Print(smallFileInvertedIndex.HashMap["organisms"].DocumentList)

		want := []string{"https://en.wikipedia.org/wiki/Anatomy", "https://en.wikipedia.org/wiki/Aquaculture",
			"https://en.wikipedia.org/wiki/Adaptive_radiation", "https://en.wikipedia.org/wiki/Abiotic_stress",
			"https://en.wikipedia.org/wiki/Apoptosis", "https://en.wikipedia.org/wiki/Asexual_reproduction",
			"https://en.wikipedia.org/wiki/Biotic", "https://en.wikipedia.org/wiki/Biochemistry",
			"https://en.wikipedia.org/wiki/Biopolymer", "https://en.wikipedia.org/wiki/Bioleaching",
			"https://en.wikipedia.org/wiki/Cell_(biology)"}
		ast.ElementsMatchf(got, want, "")
	})

	t.Run("smallGetSingleTermColor", func(t *testing.T) {
		got := smallFileInvertedIndex.getPageUrlsForTerm([]string{"color"})
		want := []string{"https://en.wikipedia.org/wiki/Amber",
			"https://en.wikipedia.org/wiki/Alpha_compositing"}
		ast.ElementsMatchf(got, want, "")
	})

	t.Run("smallGetSingleTermAbove", func(t *testing.T) {
		got := smallFileInvertedIndex.getPageUrlsForTerm([]string{"above"})
		want := []string{"https://en.wikipedia.org/wiki/Acropolis_of_Athens",
			"https://en.wikipedia.org/wiki/Adrenal_gland",
			"https://en.wikipedia.org/wiki/Afterglow"}
		ast.ElementsMatchf(got, want, "")
	})

	t.Run("smallGetSingleTermWater", func(t *testing.T) {
		got := smallFileInvertedIndex.getPageUrlsForTerm([]string{"water"})
		want := []string{"https://en.wikipedia.org/wiki/Brackish_water",
			"https://en.wikipedia.org/wiki/Extreme_poverty",
			"https://en.wikipedia.org/wiki/Transport_in_Antarctica",
			"https://en.wikipedia.org/wiki/Bubalus_(Anoa)",
			"https://en.wikipedia.org/wiki/Alkali", "https://en.wikipedia.org/wiki/Autonomous_building",
			"https://en.wikipedia.org/wiki/Beer", "https://en.wikipedia.org/wiki/Bridge",
			"https://en.wikipedia.org/wiki/Transport_in_Belgium",
			"https://en.wikipedia.org/wiki/Transport_in_Burundi",
			"https://en.wikipedia.org/wiki/Bay_(disambiguation)"}
		ast.ElementsMatchf(got, want, "")
	})

	t.Run("smallGetMissingTerm", func(t *testing.T) {
		got := smallFileInvertedIndex.getPageUrlsForTerm([]string{"supercalifragilisticexpialidocious"})
		want := []string{}
		ast.ElementsMatchf(got, want, "")
	})
}

func Test_bigGetPageUrlsForTerm(t *testing.T) {
	ast := assert.New(t)
	bigFileInvertedIndex := InvertedIndex{Filename: BigFile, HashMap: map[string]*Data{}}
	bigFileInvertedIndex.buildIndex()

	t.Run("bigGetSingleTermPineapple", func(t *testing.T) {
		got := bigFileInvertedIndex.getPageUrlsForTerm([]string{"pineapple"})
		want := []string{"https://en.wikipedia.org/wiki/Kamaka_Ukulele",
			"https://en.wikipedia.org/wiki/Hamonado",
			"https://en.wikipedia.org/wiki/Schartner_Bombe",
			"https://en.wikipedia.org/wiki/Cactus_Cooler",
			"https://en.wikipedia.org/wiki/Pineapple_coral",
			"https://en.wikipedia.org/wiki/Runts",
			"https://en.wikipedia.org/wiki/Fruit_bromelain",
			"https://en.wikipedia.org/wiki/Queens_(cocktail)",
			"https://en.wikipedia.org/wiki/George_Brown_(Australian_soccer)",
			"https://en.wikipedia.org/wiki/Chusnunia_Chalim"}
		ast.ElementsMatchf(got, want, "")
	})

	t.Run("bigGetSingleTermTrigonometry", func(t *testing.T) {
		got := bigFileInvertedIndex.getPageUrlsForTerm([]string{"trigonometry"})
		want := []string{"https://en.wikipedia.org/wiki/Standard_ruler",
			"https://en.wikipedia.org/wiki/Hyperbolic_law_of_cosines"}
		ast.ElementsMatchf(got, want, "")
	})

	t.Run("bigGetSingleTermThoracic", func(t *testing.T) {
		got := bigFileInvertedIndex.getPageUrlsForTerm([]string{"thoracic"})
		want := []string{"https://en.wikipedia.org/wiki/Subcostal_arteries",
			"https://en.wikipedia.org/wiki/Lymph_duct",
			"https://en.wikipedia.org/wiki/Articulation_of_head_of_rib",
			"https://en.wikipedia.org/wiki/American_Thoracic_Society",
			"https://en.wikipedia.org/wiki/Brachial_plexus_injury",
			"https://en.wikipedia.org/wiki/Open_aortic_surgery",
			"https://en.wikipedia.org/wiki/Lumbar_plexus",
			"https://en.wikipedia.org/wiki/Left_triangular_ligament",
			"https://en.wikipedia.org/wiki/Rhomboid_major_muscle",
			"https://en.wikipedia.org/wiki/Internal_thoracic_artery",
			"https://en.wikipedia.org/wiki/Compensatory_hyperhidrosis",
			"https://en.wikipedia.org/wiki/Boas%27_point",
			"https://en.wikipedia.org/wiki/Subclavian_steal_syndrome",
			"https://en.wikipedia.org/wiki/Chest_(journal)",
			"https://en.wikipedia.org/wiki/Association_of_Thoracic_and_Cardiovascular_Surgeons_of_Asia",
			"https://en.wikipedia.org/wiki/Thoracic_spinal_nerve_4",
			"https://en.wikipedia.org/wiki/Vukhuclepis"}
		ast.ElementsMatchf(got, want, "")
	})
}

func Test_smallSearch(t *testing.T) {
	ast := assert.New(t)
	smallFileInvertedIndex := InvertedIndex{Filename: SmallFile, HashMap: map[string]*Data{}}
	smallFileInvertedIndex.buildIndex()

	t.Run("smallSearchForApple", func(t *testing.T) {
		got := smallFileInvertedIndex.getSearchResult([]string{"apple"})
		want := []SearchResult{
			{"https://en.wikipedia.org/wiki/Apple_II_series", "Apple II series", 0.76752},
		}
		ast.ElementsMatch(got, want, "")
	})

	t.Run("smallSearchForPolitics", func(t *testing.T) {
		got := smallFileInvertedIndex.getSearchResult([]string{"Politics", "president"})
		want := []SearchResult{
			{"https://en.wikipedia.org/wiki/Brazilian_Armed_Forces", "Brazilian Armed Forces",
				1.60943},
			{"https://en.wikipedia.org/wiki/Politics_of_Armenia", "Politics of Armenia", 0.27330},
			{"https://en.wikipedia.org/wiki/Politics_of_Botswana", "Politics of Botswana", 0.23552},
			{"https://en.wikipedia.org/wiki/Politics_of_Burundi", "Politics of Burundi", 0.2299},
			{"https://en.wikipedia.org/wiki/Politics_of_Bulgaria", "Politics of Bulgaria", 0.20546},
		}
		ast.ElementsMatch(got, want, "")
	})

	t.Run("smallSearchForWaterLandAir", func(t *testing.T) {
		got := smallFileInvertedIndex.getSearchResult([]string{"water", "land", "air"})
		want := []SearchResult{
			{"https://en.wikipedia.org/wiki/Air_show", "Air show", 0.72973},
			{"https://en.wikipedia.org/wiki/Bay_(disambiguation)", "Bay (disambiguation)",
				0.7092499999999999},
			{"https://en.wikipedia.org/wiki/Arable_land", "Arable land", 0.58881},
			{"https://en.wikipedia.org/wiki/Brackish_water", "Brackish water", 0.46977},
			{"https://en.wikipedia.org/wiki/Transport_in_Belgium", "Transport in Belgium",
				0.41166},
		}
		ast.ElementsMatch(got, want, "")
	})

}

func Test_bigSearch(t *testing.T) {
	ast := assert.New(t)
	bigFileInvertedIndex := InvertedIndex{Filename: BigFile, HashMap: map[string]*Data{}}
	bigFileInvertedIndex.buildIndex()

	t.Run("bigSearchForWatermelonPotato", func(t *testing.T) {
		got := bigFileInvertedIndex.getSearchResult([]string{"watermelon", "potato"})
		want := []SearchResult{
			{"https://en.wikipedia.org/wiki/Melody_potato", "Melody potato", 1.41344},
			{"https://en.wikipedia.org/wiki/Idaho_Potato_Museum", "Idaho Potato Museum",
				1.06008},
			{"https://en.wikipedia.org/wiki/Russian_blue_potato", "Russian blue potato",
				0.99772},
			{"https://en.wikipedia.org/wiki/Idaho_Potato_Commission",
				"Idaho Potato Commission", 0.94229},
			{"https://en.wikipedia.org/wiki/One_potato,_two_potato",
				"One potato, two potato", 0.90864},
		}
		ast.ElementsMatch(got, want, "")
	})

	t.Run("bigSearchForScrewdriver", func(t *testing.T) {
		got := bigFileInvertedIndex.getSearchResult([]string{"screwdriver"})
		want := []SearchResult{
			{"https://en.wikipedia.org/wiki/Cam_out", "Cam out", 0.44899},
			{"https://en.wikipedia.org/wiki/Staple_knocker", "Staple knocker", 0.29746},
		}
		ast.ElementsMatch(got, want, "")
	})

	t.Run("bigSearchForPolitics", func(t *testing.T) {
		got := bigFileInvertedIndex.getSearchResult([]string{"politics"})
		want := []SearchResult{
			{"https://en.wikipedia.org/wiki/John_Murray_(colonial_administrator)",
				"John Murray (colonial administrator)", 1.31803},
			{"https://en.wikipedia.org/wiki/Office_Politics", "Office Politics",
				1.31803},
			{"https://en.wikipedia.org/wiki/Politics_of_Washington",
				"Politics of Washington", 1.09835},
			{"https://en.wikipedia.org/wiki/Politics_(1940s_magazine)",
				"Politics (1940s magazine)", 0.77531},
			{"https://en.wikipedia.org/wiki/Faith_and_Politics_Institute",
				"Faith and Politics Institute", 0.73223},
		}
		ast.ElementsMatch(got, want, "")
	})

	t.Run("bigSearchForBananaCheesecake", func(t *testing.T) {
		got := bigFileInvertedIndex.getSearchResult([]string{"Banana", "Cheesecake"})
		want := []SearchResult{
			{"https://en.wikipedia.org/wiki/Banana_Ridge", "Banana Ridge", 1.75258},
			{"https://en.wikipedia.org/wiki/Masak_Hijau_banana", "Masak Hijau banana",
				0.87629},
			{"https://en.wikipedia.org/wiki/Musa_thomsonii", "Musa thomsonii", 0.87629},
			{"https://en.wikipedia.org/wiki/Banana_Xanthomonas_wilt",
				"Banana Xanthomonas wilt", 0.79663},
			{"https://en.wikipedia.org/wiki/Cane_knife", "Cane knife", 0.79663},
		}
		ast.ElementsMatch(got, want, "")
	})

	t.Run("bigSearchForNonExistingTerm", func(t *testing.T) {
		got := bigFileInvertedIndex.getSearchResult([]string{"pterodactyl"})
		want := []SearchResult{}
		ast.ElementsMatch(got, want, "")
	})

}
