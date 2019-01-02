package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

type Validation struct {
	Key       string `yaml:"Key"`
	Validated string `yaml:"Validated"`
	Learning  string `yaml:"Learning"`
	Link      string `yaml;"Link"`
}

type Validations []Validation

func (v Validations) HTML() (t string) {
	t += "<ul>"
	for _, i := range v {

		color := "green"

		switch i.Validated {
		case "-1":
			color = "red"
		case "1":
			color = "green"
		default:
			color = "orange"
		}
		t += fmt.Sprintf(`<li style="color: %s">%s</li>`, color, i.Key)
	}

	t += "</ul>"

	return t
}

type Canvas struct {
	Name string `yaml:"Name"`

	Problem struct {
		Validations          Validations `yaml:"Validations"`
		ExistingAlternatives Validations `yaml:"ExistingAlternatives"`
	} `yaml:"Problem"`

	Solution   Validations `yaml:"Solution"`
	KeyMetrics Validations `yaml:"KeyMetrics"`

	UniqueValueProposition struct {
		Validations      Validations `yaml:"Validations"`
		HighLevelConcept Validations `yaml:"HighLevelConcept"`
	} `yaml:"UniqueValueProposition"`

	UnfairAdvantage Validations `yaml:"UnfairAdvantage"`
	Channels        Validations `yaml:"Channels"`

	CustomerSegments struct {
		Validations   Validations `yaml:"Validations"`
		EarlyAdopters Validations `yaml:"EarlyAdopters"`
	} `yaml:"CustomerSegments"`

	CostStructure  Validations `yaml:"CostStructure"`
	RevenueStreams Validations `yaml:"RevenueStreams"`
}

func (c Canvas) HTML() string {
	return fmt.Sprintf(`
	<html>
<head>
	<title>GitCanvas</title>

	<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.1.3/css/bootstrap.min.css" integrity="sha384-MCw98/SFnGE8fJT3GXwEOngsV7Zt27NXFoaoApmYm81iuXoPkFOJwJ8ERdknLPMO" crossorigin="anonymous">
</head>


<body>
    <!-- Google Webfonts -->

	<style>
	main {
		margin: 20px 0 20px 0;
		font-family: "roboto";
	}

	.column {
		border: 1px solid #AAAAAA;
		overflow: hidden;
	}

	.row-one .column {
		min-height: 300px;
	}

	.row-two .column {
		min-height: 150px;
	}

	hr.divider {
		margin: 0 -20px;
		border-top: 1px solid #AAAAAA
	}

	/* Typography */
	h2.section-title {
		font-size: 17px;
		font-weight: 500;
		color: #333333;
	}

	h3.sub-section-title {
		font-size: 14px;
		font-weight: 500;
	}

	section p {
		color: #666666;
		font-weight: 300;
	}

	/* Equal Height Columns */

	.row-eq-height {
		display: -webkit-box;
		display: -webkit-flex;
		display: -ms-flexbox;
		display: flex;
	}
</style>

	<main>
        <div class="container-fluid">
            <!-- start Bootstrap Grid -->
            <div class="row row-header section-container">
                <div class="col-md-10 col-md-offset-1">
                    <h1>%s</h1>


                </div>
            </div>

            <!-- Lean Canvas is 2 rows, with 5 columns. The bootstrap grid only contains 12 columns, which isn't evenly divisble by 5, so we will instead use 10 columns only, offsetting 2 of the columns. -->
            <div class="row row-one section-container row-eq-height">
                <!-- begin first row -->
                <section class="column col-md-2 col-md-offset-1" id="problem">
                    <h2 class="section-title">Problem</h2>
					%s
					<h3 class="sub-section-title">Existing Alternatives</h3>
					%s
                </section>

                <div class="column col-md-2">
                    <section class="top" id="solution">
                        <h2 class="section-title">Solution</h2>
                        %s
                    </section>
                    <hr class="divider">
                    </hr>
                    <section class="bottom" id="key-metrics">
                        <h2 class="section-title">Key Metrics</h2>
                        %s
                    </section>
                </div>

                <section class="column col-md-2" id="unique-value-proposition">
                    <h2 class="section-title">Unique Value Proposition</h2>
					%s

					<h3 class="sub-section-title">High-Level Concept</h3>
					%s
                </section>

                <div class="column col-md-2">
                    <section class="" id="unfair-advantage">
                        <h2 class="section-title">Unfair Advantage</h2>
						%s
                    </section>
                    <hr class="divider">
                    </hr>
                    <section class="bottom" id="channels">
                        <h2 class="section-title">Channels</h2>
						%s
                    </section>
                </div>

                <section class="column col-md-2" id="customer-segments">
                    <h2 class="section-title">Customer Segments</h2>
					%s

					<h3 class="sub-section-title">Early Adopters</h3>
					%s
                </section>

            </div><!-- End first row -->

            <div class="row row-two section-container row-eq-height">
                <!-- Begin second row -->
                <section class="column col-md-5 col-md-offset-1" id="cost-structure">
                    <h2 class="section-title">Cost Structure</h2>
					%s
                </section>
                <section class="column col-md-5" id="revenue-streams">
                    <h2 class="section-title">Revenue Streams</h2>
                    %s
                </section>
            </div> <!-- End second row -->
        </div>
    </main>
</body>
</html>`, c.Name, c.Problem.Validations.HTML(), c.Problem.ExistingAlternatives.HTML(), c.Solution.HTML(), c.KeyMetrics.HTML(), c.UniqueValueProposition.Validations.HTML(), c.UniqueValueProposition.HighLevelConcept.HTML(),
		c.UnfairAdvantage.HTML(), c.Channels.HTML(), c.CustomerSegments.Validations.HTML(), c.CustomerSegments.EarlyAdopters.HTML(), c.CostStructure.HTML(), c.RevenueStreams.HTML())
}

func main() {
	flag.Parse()
	fileName := flag.Arg(0)

	canvas := &Canvas{}

	b, _ := ioutil.ReadFile(fileName)

	err := yaml.Unmarshal(b, canvas)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(canvas.HTML())
}
