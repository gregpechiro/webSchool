var fileType;
var current;
var settings;
var editor;
var tree;
var memFiles = {};
$(document).ready(function() {
    settings = getSettings();

    editor = ace.edit("editor");
    
    var snippetManager = ace.require("ace/snippets").snippetManager;

    function onKeyDown(e) {
        // ctrl+r remove default
        if (e.ctrlKey) { // ctrl
            if (e.keyCode == 83) { // +s
                e.preventDefault();
                if (e.shiftKey) { // + shift
                    saveAll();
                    return;
                }
                save(current);
                return;
            }

            if (e.keyCode == 73) { // +i
                e.preventDefault();
                newImage(current);
                return;
            }

            if (e.altKey) { // + alt
                e.preventDefault();
                if (e.keyCode == 78) { // + n
                    if (e.shiftKey) { // + shift
                        newFolder(current);
                        return;
                    }
                    newFile(current);
                    return;
                }
            }
        }
    }
    // register the handler
    document.addEventListener('keydown', onKeyDown, false);

    $('button#save').click(function() {
        save(current);
    });

    $('button#saveAll').click(function() {
        saveAll();
    });

    $('a#newFile').click(function() {
        newFile(current);
    });

    $('a#newFolder').click(function() {
        newFolder(current);
    });

    $('a#newImage').click(function() {
        newImage(current);
    });

    $(function () {
        $('[data-toggle="popover"]').popover();
    });

    $('#newFileModal').on('hidden.bs.modal', function (e) {
        $('form#newFileForm')[0].reset();
    });

    $('#newFolderModal').on('hidden.bs.modal', function (e) {
        $('form#newFolderForm')[0].reset();
    });

    $('form#newFileForm').submit(function(e) {
        e.preventDefault();
        newFileSubmit();
    });

    $('button#newFile').click(function() {
        newFileSubmit();
    });

    $('form#newFolderForm').submit(function(e) {
        e.preventDefault();
        newFolderSubmit();
    });

    $('button#newFolder').click(function() {
        newFolderSubmit();
    });

    $('button#upload').click(function() {
        var formData = new FormData($('#uploader')[0]);
        $.ajax({
            url: '/project/' + project + '/upload',
            method: 'POST',
            cache: false,
            contentType: false,
            processData: false,
            data: formData,
            success: function(resp) {
                // check for returned error
                if (resp.error) {
                    displayError(resp.output);
                    return
                }
                tree.refresh();
                displaySuccess(resp.output);
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
        $('#uploadModal').modal('hide');
    });

});

$('html, body').height(window.innerHeight - 52);

function getFile(id) {
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
        var memFile = memFiles[id];
        if (memFile != null) {
            if (memFile.unsaved) {
                memFile.save();
                return;
            }
        }
        displayError('No changes detected')
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

function newImage(id) {
    if (id != undefined && id != '') {
        var n = tree.get_node(id);
        while (n.type !== 'dir' && n.id !== '#') {
            n = tree.get_node(n.parent);
        }
        $('input#uploadPath').val(n.id);
        var p = '';
        if (n.id != '#') {
            p = decodeURIComponent(n.id);
            if (p[0] == '/') {
                p = p.slice(1);
            }
            if (p[p.length - 1] != '/') {
                p += '/';
            }
        }
        $('label#uploadPath').text(p);
    } else {
        $('input#uploadPath').val('#');
    }
    $('div#uploadModal').modal('show');
}

function newFile(id) {
    if (id != undefined && id != '') {
        var n = tree.get_node(id);
        while (n.type !== 'dir' && n.id !== '#') {
            n = tree.get_node(n.parent);
        }
        $('input#filePath').val(n.id);
        var p = '';
        if (n.id != '#') {
            p = decodeURIComponent(n.id);
            if (p[0] == '/') {
                p = p.slice(1);
            }
            if (p[p.length - 1] != '/') {
                p += '/';
            }
        }
        $('label#filePath').text(p);
    } else {
        $('input#filePath').val('#');
    }
    $('div#newFileModal').modal('show');
}

function newFileSubmit() {
    $.ajax({
        url: '/project/' + project + '/addFile',
        method: 'POST',
        data: formToObject($('#newFileForm')),
        success: function(resp) {
            // check for returned error
            if (resp.error) {
                displayError(resp.output);
                return
            }
            tree.refresh();
            displaySuccess(resp.output);
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
    $('#newFileModal').modal('hide');
}

function newFolder(id) {
    if (id != undefined && id != '') {
        var n = tree.get_node(id);
        while (n.type !== 'dir' && n.id !== '#') {
            n = tree.get_node(n.parent);
        }
        $('input#folderPath').val(n.id);
        var p = '';
        if (n.id != '#') {
            p = decodeURIComponent(n.id);
            if (p[0] == '/') {
                p = p.slice(1);
            }
            if (p[p.length - 1] != '/') {
                p += '/';
            }
        }
        $('label#folderPath').text(p);
    } else {
        $('input#folderPath').val('#');
    }
    $('div#newFolderModal').modal('show');
}

function newFolderSubmit() {
    $.ajax({
        url: '/project/' + project + '/mkdir',
        method: 'POST',
        data: formToObject($('#newFolderForm')),
        success: function(resp) {
            // check for returned error
            if (resp.error) {
                displayError(resp.output);
                return
            }
            tree.refresh();
            displaySuccess(resp.output);
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
    $('#newFolderModal').modal('hide');
}

function setEditorHeader(header) {
    $('div#fileName').text(header);
}
