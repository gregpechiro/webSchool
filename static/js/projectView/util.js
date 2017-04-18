// save settings to local storage
function saveSettings(settings) {
    var s = JSON.stringify(settings);
    s = btoa(unescape(encodeURIComponent(s)));
    localStorage.setItem('settings', s);
}

// return saved settings from local storage
function getSettings() {
    var s = localStorage.getItem('settings');
    if (s !== undefined && s !== null && s !== '') {
        s = decodeURIComponent(escape(atob(s)));
        if (s !== '' && s !== undefined && s[0] === '{') {
            return JSON.parse(s);
        }
        saveSettings({'editor':{}})
    }
    return {'editor':{}};
}
