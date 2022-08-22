module getHackernewsData

go 1.19

require (
	github.com/aws/aws-lambda-go v1.34.1
	github.com/aws/aws-sdk-go v1.35.33
	pht/hndata v0.0.0-00010101000000-000000000000
)

replace pht/hndata => ./hndata
