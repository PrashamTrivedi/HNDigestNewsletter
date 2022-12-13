module getHackernewsData

go 1.19

require (
	github.com/aws/aws-lambda-go v1.36.0
	github.com/aws/aws-sdk-go v1.44.159
	pht/hndata v0.0.0-00010101000000-000000000000
)

require github.com/jmespath/go-jmespath v0.4.0 // indirect

replace pht/hndata => ./hndata
