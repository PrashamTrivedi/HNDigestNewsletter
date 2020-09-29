import {getFirebaseData} from './getFirebaseData'

export const getFBData = async ()=>{
    await getFirebaseData({type:'job'})
}

module.exports.getFBData()