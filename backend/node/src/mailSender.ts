import AWS from 'aws-sdk'
AWS.config.region = 'ap-south-1'
const ses = new AWS.SES()

export const sendEmail = async (event: any) => {
    const body = JSON.parse(event.Records[0].body)
    const html = body.html
    const type = body.type
    const toEmail = body.to || ['prash2488@gmail.com']
    const subjectLine = `Your Hackernews Digest for ${type}`
    const fromEmail = 'jobs@prashamhtrivedi.in'
    const fromBase64 = Buffer.from(fromEmail).toString('base64')
    const toAddresses = toEmail

    const htmlBody = html

    const emailParams = {
        Destination: {
            ToAddresses: toAddresses,
        },
        Message: {
            Body: {
                Html: {
                    Charset: 'UTF-8',
                    Data: htmlBody,
                },
            },
            Subject: {
                Charset: 'UTF-8',
                Data: subjectLine,
            },
        },
        ReplyToAddresses: [fromEmail],
        Source: `=?utf-8?B?${fromBase64}?= <${fromEmail}>`,
    }

    const sesResponse = await ses.sendEmail(emailParams).promise()

    return {
        statusCode: 200
    }
}

