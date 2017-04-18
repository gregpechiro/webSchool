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

var fileType;
var path;

var settings = getSettings();

var editor = ace.edit("editor");
var t;

function getFile(id) {
    $.ajax({
        url: '/project/' + project + '/file?path=' + id,
        method: 'GET',
        success: function(resp) {
            // check for returned error
            if (resp.error) {
                alert('error');
                return
            }
            // replace editor with formated code
            var file = atob(resp.output);
            f = resp.output
            editor.setValue(file, 1);
            editor.session.setMode("ace/mode/" + fileType);
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
