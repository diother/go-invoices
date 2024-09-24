package documents

const (
	marginTop    = 32
	marginLeft   = 40
	marginRight  = 555
	marginBottom = 810
)

const (
	itemHeight             = 50
	firstPageStartY        = 357.0
	secondPageStartY       = 135.0
	firstPageCapacity      = 8
	subsequentPageCapacity = 12
	firstPageTableY        = 315
	subsequentPageTableY   = 93
)

type DocumentService struct{}

func NewDocumentService() *DocumentService {
	return &DocumentService{}
}
