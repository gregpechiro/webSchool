function MemFile(id) {
    this.id = id;
}

MemFile.prototype = {
    id: '',
    fresh: true,
    unsaved: false,
    update: function() {
        if (this.unsaved) {
            return;
        }
        this.unsaved = true;
        var node = tree.get_node(this.id);
        node.text = '* ' + node.text;
        tree.redraw(node);
    },
    save: function() {
        var memFile = this;
        $.ajax({
            url: '/project/' + project + '/file/save',
            method: 'POST',
            data: {
                path: this.id,
                data: memFile.session.getValue()
            },
            success: function(resp) {
                // check for returned error
                if (resp.error) {
                    $.Notification.autoHideNotify('error', 'top center', resp.output);
                    return
                }
                memFile.unsaved = false;
                var node = tree.get_node(memFile.id);
                node.text = node.text.replace('* ', '');
                tree.redraw(node);
                $.Notification.autoHideNotify('success', 'top center', resp.output);
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
};
