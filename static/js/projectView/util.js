// save settings to local storage
function saveSettings(settings) {
    var s = JSON.stringify(settings);
    setLocal('settings', s);
}

// return saved settings from local storage
function getSettings() {
    var s = getLocal('settings');
    if (s !== '' && s !== undefined && s[0] === '{') {
        return JSON.parse(s);
    }
    saveSettings({'editor':{}});
    return {'editor':{}};
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

function formToObject(form) {
	var object = {};
	var formArray = form.serializeArray();
	$.each(formArray, function() {
		if (object[this.name] !== undefined) {
			if (!object[this.name].push) {
				object[this.name] = [object[this.name]];
			}
	    		object[this.name].push(this.value || '');
		} else {
			object[this.name] = this.value || '';
		}
	});
	return object;
};
