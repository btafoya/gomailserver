package totp

import (
	"github.com/pquerna/otp/totp"
)

type TOTPService struct {
	issuer string
}

func NewTOTPService(issuer string) *TOTPService {
	return &TOTPService{issuer: issuer}
}

func (s *TOTPService) GenerateSecret(email string) (*TOTPSetup, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      s.issuer,
		AccountName: email,
	})
	if err != nil {
		return nil, err
	}

	return &TOTPSetup{
		Secret:     key.Secret(),
		URL:        key.URL(),
		QRCodeData: generateQRCode(key.URL()),
	}, nil
}

func (s *TOTPService) Validate(secret, code string) bool {
	return totp.Validate(code, secret)
}

type TOTPSetup struct {
	Secret     string
	URL        string
	QRCodeData []byte // PNG image data
}

func generateQRCode(url string) []byte {
	// Placeholder for QR code generation logic.
	// This would typically use a library like github.com/skip2/go-qrcode
	return nil
}
