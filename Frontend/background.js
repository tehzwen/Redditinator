// listen for icon to be clicked
console.log("background running");
chrome.browserAction.onClicked.addListener(buttonClicked)

function buttonClicked(tab){
    console.log(tab);
}