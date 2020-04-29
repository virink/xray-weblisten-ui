module.exports = {
    publicPath: process.env.NODE_ENV === 'production' ?
        '/static/' : '/',
    chainWebpack: config => {
        config
            .plugin('html')
            .tap(args => {
                args[0].title = 'Xray 被动扫描管理'
                return args
            })
    }
}