package searchengine

import (
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
