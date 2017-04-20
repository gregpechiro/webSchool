var fileType;
var current;
var settings;
var editor;
var tree;
var memFiles = {};

$(document).ready(function() {
    settings = getSettings();
    editor = ace.edit("editor");

    function onKeyDown(e) {
        // ctrl+r remove default
        if (e.ctrlKey) { // ctrl
            if (e.keyCode == 83) { // +s
                e.preventDefault();
                if (e.shiftKey) {
                    saveAll();
                    return;
                }
                save(current);
                return
            }
        }
    }

    // register the handler
    document.addEventListener('keydown', onKeyDown, false);
});

$('html, body').height(window.innerHeight - 52);

$('.away').click(function() {
    parent = $(this).closest('.side');
    if (parent.width() > 35) {
        parent.find('.content').addClass('hide');
        parent.find('.alt').removeClass('hide');
        parent.width(35);
        return
    }
    parent.find('.alt').addClass('hide');
    parent.find('.content').removeClass('hide');
    parent.width('15%');
});


function getFile(id) {
    if (id == current) {
        return;
    }
    current = id;

    // check cache for file before requesting from server
    var memFile = memFiles[id];
    if (memFile != null && memFile.id != '') {
        editor.setSession(memFile.session)
        return;
    }


    $.ajax({
        url: '/project/' + project + '/file?path=' + id,
        method: 'GET',
        success: function(resp) {
            // check for returned error
            if (resp.error) {
                alert('error');
                return
            }
            // parse returned file
            var file = atob(resp.output);

            // create and set new chached file
            memFile = new MemFile(id);
            memFile.session = ace.createEditSession(file, "ace/mode/" + resp.fileType);
            memFiles[id] = memFile;

            // replace editor with formated code
            editor.setSession(memFile.session);

            return
        },
        // display server error
        error: function(e, d) {
            alert('error');
            console.log(e);
            console.log(d);
            return
        }
    });
}

function save(id) {
    if (id != '') {
        var memFile = memFiles[current];
        if (memFile != null) {
            if (memFile.unsaved) {
                memFile.save();
                return;
            }
        }
        $.Notification.autoHideNotify('error', 'top center', 'No changes detected');
    }
}

function saveAll() {
    for (var key in memFiles) {
        if (!memFiles.hasOwnProperty(key)) {
            continue;
        }
        var memFile = memFiles[key];
        if (memFile != null) {
            if (memFile.unsaved) {
                memFile.save();
            }
        }
    }
}

function isSaved(id) {
    var memFile = memFiles[id];
    if (memFile != null) {
        if (memFile.unsaved) {
            return false;
        }
    }
    return true;
}
