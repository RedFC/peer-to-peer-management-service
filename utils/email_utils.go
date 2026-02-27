package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/skip2/go-qrcode"
)

// func EmailBodyTemplate(fName string, lName string, qrCodeImg string, ios_url string, android_url string) string {

// 	var image = "<img src=\"" + qrCodeImg + "\" />"

// 	return fmt.Sprintf(`
// 		Hi %s %s,

// 		Welcome to Peer-Peer Communication platform!

// 		Download the App:
// 		Use the links below to install the Peer-Peer Communication app on your device:
// 		iOS: %s
// 		Android: %s

// 		Activate Your Access:
// 		After installation, open the app and scan the below QR code to complete your onboarding and activation.

// 		%s

// 		Thanks,
// 		Team Peer-to-Peer Communication`,
// 		fName,
// 		lName,
// 		ios_url,
// 		android_url,
// 		image,
// 	)
// }

func EmailBodyTemplate(
	fName string,
	lName string,
	qrCodeImg string,
	iosURL string,
	androidURL string,
) string {

	if iosURL == "" {
		iosURL = "https://apps.apple.com/us/app/peer-peer-communication/id6441234567"
	}

	if androidURL == "" {
		androidURL = "https://play.google.com/store/apps/details?id=com.peerpeercommunication"
	}

	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Peer-to-Peer Communication Invitation</title>
</head>
<body style="font-family: Arial, Helvetica, sans-serif; background-color: #f5f5f5; padding: 20px;">

	<table width="100%" cellpadding="0" cellspacing="0">
		<tr>
			<td align="center">
				<table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff; padding: 24px; border-radius: 6px;">
					
					<tr>
						<td>
							<p>Hi <strong>%s %s</strong>,</p>

							<p>
								Welcome to the <strong>Peer-to-Peer Communication</strong> platform.
							</p>

							<p>
								<strong>Download the app:</strong><br>
								<a href="%s">iOS App Store</a><br>
								<a href="%s">Google Play Store</a>
							</p>

							<p>
								<strong>Activate your access:</strong><br>
								After installing the app, open it and scan the QR code below to complete your onboarding.
							</p>

							<p style="text-align: center; margin: 20px 0;">
								<img src="%s" alt="Activation QR Code" style="max-width: 200px; height: auto; border: 1px solid #ddd; padding: 5px; background: #fff;" />
							</p>

							<p>
								If you have any questions, please contact your administrator.
							</p>

							<p>
								Regards,<br>
								<strong>Team Peer-to-Peer Communication</strong>
							</p>
						</td>
					</tr>

				</table>
			</td>
		</tr>
	</table>

</body>
</html>
`,
		fName,
		lName,
		iosURL,
		androidURL,
		qrCodeImg,
	)
}

func GenerateQRCode(payload interface{}) (string, error) {
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	png, err := qrcode.Encode(string(jsonBytes), qrcode.Medium, 256)
	if err != nil {
		return "", err
	}

	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(png), nil
}
