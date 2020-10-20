import AWS from 'aws-sdk'
import crypto from 'crypto'
const dynamoDb = new AWS.DynamoDB()
const documentClient = new AWS.DynamoDB.DocumentClient()
const params = {
    TableName: 'hndigest_emails',
    Item: {}
}
// Difining algorithm 
const algorithm = 'aes-256-cbc'

// Defining key 
const key = crypto.randomBytes(32)

// Defining iv 
const iv = crypto.randomBytes(16)

export const registerEmailLambda = async (event: any) => {
    const emailId = event.email
    const types = event.types



    const {encryptedData, ivData} = encrypt(`${JSON.stringify(event)}`)
    const decryptData = decrypt(encryptedData, ivData).decryptedData

    params.Item = {
        emailId: emailId,
        types: types,
        isRegistered: false,
        token: encryptedData,
        decryptedData: JSON.parse(decryptData)
    }

    console.log(params.Item)
    // const data = await documentClient.put(params).promise()


    return {
        statusCode: 200,
        message: 'Welcome, We have sent you the verification email.',
        params
    }
}

function encrypt(text: string) {

    // Creating Cipheriv with its parameter 
    let cipher = crypto.createCipheriv('aes-256-cbc', Buffer.from(key), iv)

    // Updating text 
    let encrypted = cipher.update(text)

    // Using concatenation 
    encrypted = Buffer.concat([encrypted, cipher.final()])

    // Returning iv and encrypted data 
    return {
        ivData: iv.toString('hex'),
        encryptedData: encrypted.toString('hex')
    }
}

function decrypt(encryptedText: string, ivData: string) {
    const ivBuffer = Buffer.from(ivData, 'hex')
    const encryptedBuffer = Buffer.from(encryptedText, 'hex')
    let cipher = crypto.createDecipheriv('aes-256-cbc', Buffer.from(key), ivBuffer)
    let decrypted = cipher.update(encryptedBuffer)

    // Using concatenation 
    decrypted = Buffer.concat([decrypted, cipher.final()])

    // Returning iv and encrypted data 
    return {
        decryptedData: decrypted.toString()
    }
}