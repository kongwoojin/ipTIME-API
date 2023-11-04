package routers

// Desktop UI
const routerLogin = "login.cgi"
const routerRoot = "timepro.cgi"
const routerLoginHandler = "login_handler.cgi"
const routerLoginSession = "login_session.cgi"
const routerSystemInfoStatus = "timepro.cgi?tmenu=iframe&smenu=system_info_status"
const routerLanPCInfoStatus = "timepro.cgi?tmenu=iframe&smenu=lan_pcinfo_status"
const routerPortForwardList = "timepro.cgi?tmenu=iframe&smenu=user_portforward&mode=user"
const routerPortForwardRulesDownload = "download_portforward.cgi"
const routerPortForwardRestore = "timepro.cgi?tmenu=iframe&smenu=restore_portforward"

const routerWOLList = "timepro.cgi?tmenu=iframe&smenu=expertconfwollist"

const routerMacManagementModeSubmit = "timepro.cgi?tmenu=iframe&smenu=macauth_bsslist_submit"
const routerMacManagementModeList = "timepro.cgi?tmenu=iframe&smenu=macauth_bsslist"

// Mobile UI
const mobileRouterMacAuth = "wirelessconf/macauth/iux.cgi"
const mobileRouterSubmit = "cgi/iux_set.cgi"
