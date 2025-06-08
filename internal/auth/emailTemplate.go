package auth

import "fmt"

func generateTemplate(otp string) string {
	return fmt.Sprintf(emailTemplate, styles, otp)
}

var styles string = `
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #111111;
            margin: 0;
            padding: 0;
        }
        .container {
            width: 100vw;
            max-width: 800px;
            margin: 20px auto;
            background: #111111;
            color: #eee;
            border-radius: 8px;
            box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.1);
            padding: 20px;
            text-align: center;
        }
        .header {
            font-size: 28px;
            font-weight: bold;
            color: #eee;
        }
        .otp-box {
            font-size: 28px;
            font-weight: bold;
            color: #ffffff;
            background: #F0841E;
            display: inline-block;
            padding: 10px 20px;
            margin: 20px 0;
            border-radius: 5px;
        }
        .text {
            font-size: 20px;
        }
        .footer {
            font-size: 16px;
            color: #bbb;
            margin-top: 20px;
        }
    </style>
    `

var emailTemplate = `
<!DOCTYPE html>
<html>
<head>
    %s
</head>
<body>
    <div class="container">
        <p class="header">Account Verification</p>
        <p class="text">Your One Time Password (OTP) is:</p>
        <p class="otp-box">%s</p>
        <p class="text">Please enter this code to verify your account. This OTP will expire in 5 minutes.</p>
        <p class="footer">If you didn’t request this, you can safely ignore this email.</p>
    </div>
</body>
</html>
`
