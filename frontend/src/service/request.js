import axios from 'axios'
// import {
//     Message
// } from 'element-ui'

import baseURL from '../config'

const service = axios.create({
    baseURL: baseURL,
    withCredentials: false,
    timeout: 5000
})

export default service