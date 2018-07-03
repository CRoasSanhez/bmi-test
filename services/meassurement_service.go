package services
import(
	"bodyMaxIndex/repositories"
)
// MeassurementService ...
type MeassurementService interface{}

// meassurementService ...
type meassurementService struct {
	repo repositories.UserRepository
}