var browser = browser || chrome
console.log(browser)
browser.extension.onConnect.addListener(function (port) {
    port.onMessage.addListener(function (msg) {
        console.log(msg);
    });
})