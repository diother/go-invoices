package documents

const (
	marginTop    = 32
	marginLeft   = 40
	marginRight  = 555
	marginBottom = 810
)

type DocumentService struct{}

func NewDocumentService() *DocumentService {
	return &DocumentService{}
}
