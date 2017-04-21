// save settings to local storage
function saveSettings(settings) {
    var s = JSON.stringify(settings);
    setLocal('settings', s);

    /*s = btoa(unescape(encodeURIComponent(s)));
    localStorage.setItem('settings', s);*/
}

// return saved settings from local storage
function getSettings() {
    var s = getLocal('settings');
    if (s !== '' && s !== undefined && s[0] === '{') {
        return JSON.parse(s);
    }
    saveSettings({'editor':{}});
    return {'editor':{}};

    /*var s = localStorage.getItem('settings');
    if (s !== undefined && s !== null && s !== '') {
        s = decodeURIComponent(escape(atob(s)));
        if (s !== '' && s !== undefined && s[0] === '{') {
            return JSON.parse(s);
        }
        saveSettings({'editor':{}})
    }
    return {'editor':{}};*/
}

function setLocal(key, data) {
    if (key == '' || data == '') {
        return false
    }
    data = btoa(unescape(encodeURIComponent(data)));
    localStorage.setItem(key, data);
}

function getLocal(key) {
    var s = localStorage.getItem(key);
    if (s !== undefined && s !== null && s !== '') {
        return decodeURIComponent(escape(atob(s)));
    }
    return '';
}

function displayError(msg) {
    $.Notification.autoHideNotify('error', 'top center', msg);
}

function displaySuccess(msg) {
    $.Notification.autoHideNotify('success', 'top center', msg);
}
