import {getFirebaseData} from './getFirebaseData'
import {registerEmailLambda} from './registerEmailLambda'


export const getFBData = async () => {
    await getFirebaseData({type:'show'})
    // await registerEmailLambda({email: 'prash2488@gmail.com', types: ['show', 'ask']})
}

module.exports.getFBData()