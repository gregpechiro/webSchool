$(document).ready(function() {
    $('#filetree').on('select_node.jstree', function(e, data) {

        var evt = window.event || e;
        var button = evt.which || evt.button;
        if ( button != 1 ) {
            return;
        }

        var n = data.node;
        if (n.type === 'dir') {
            if (t.is_closed(n)) {
                t.open_node(n);
            } else {
                t.close_node(n);
            }
            return;
        }
        fileType = n.type;
        path = n.id;
        getFile(n.id);

    }).on('move_node.jstree', function(e, data) {
        var frm = data.node.id;
        // var to  = data.parent + '/' + data.node.text;
        var to  = ((data.parent == "#") ? '/' + data.node.text : data.parent + '/' + data.node.text);
        if (frm !== to) {
            var form = $('<form method="post" class="hide" action="/project/' + project + '/file/move"><input name="to" value="' + to + '"><input name="from" value="' + frm + '"><input name="type" value="mov"></form>')
            $('body').append(form);
            form.submit();
        }
    }).jstree({
        "core" : {
            "multiple": false,
            "animation" : 0,
            "check_callback" : true,
            "themes" : {
                "stripes" : true
            },
            "data" : {
                "url" : function(node) {
                    if (node.id == '#') {
                        return '/project/' + project + '/files';
                    }
                    return '/project/' + project + '/files?path=' + node.id;
                },
                "data" : function (node) {
                    return node.id == "#" ? {} : { "id" : node.id }
                }
            }
        },
        "sort": function(n1, n2) {
            if (this.get_type(n1) != this.get_type(n2)) {
                if (this.get_type(n1) == "dir") {
                    return -1
                }
                return 1
            }
            if (this.get_text(n1) <= this.get_text(n2)) {
                return -1
            }
            return 1
        },
        "types": {
            "#": {
                "max_children" : 1,
                "valid_children" : ["dir", "file"]
            },
            "dir" : {
                "icon" : "glyphicon glyphicon-folder-open",
                "valid_children" : ["dir", "file"]
            },
            "html" : {
                "icon" : "fa fa-file-code-o",
                "valid_children" : []
            },
            "css" : {
                "icon" : "fa fa-file-code-o",
                "valid_children" : []
            },
            "javascript" : {
                "icon" : "fa fa-file-code-o",
                "valid_children" : []
            }
        },
        "contextmenu" : {
            items : {
                "new" : {
                    "separator_before"  : false,
                    "separator_after"   : true,
                    "label"             : "New",
                    "action"            : false,
                    "submenu" :{
                        "create_file" : {
                            "label" : "File",
                            action : function (obj) {
                                var n = t.get_node(obj.reference[0].id);
                                while (n.type !== 'dir' && n.id !== '#') {
                                    n = t.get_node(n.parent);
                                }
                                $('input#filePath').val(n.id);
                                $('div#newFileModal').modal('show');
                            }
                        },
                        "create_folder" : {
                            "seperator_before" : false,
                            "seperator_after" : false,
                            "label" : "Folder",
                            action : function (obj) {
                                var n = t.get_node(obj.reference[0].id);
                                while (n.type !== 'dir' && n.id !== '#') {
                                    n = t.get_node(n.parent);
                                }
                                $('input#folderPath').val(n.id);
                                $('div#newFolderModal').modal('show');
                            }
                        }
                    }
                },
                "rename": {
                    "separator_before"  : false,
                    "separator_after"   : false,
                    "label"             : "Rename",
                    "action"            : function(obj) {
                        var n = t.get_node(obj.reference[0].id);
                        var old_name = n.text;
                        t.edit(n, 0, function(node, status, cancel) {
                            if (old_name == node.text) {
                                return
                            }
                            var frm = node.id;
                            var to  = ((node.parent == "#") ? '/' + node.text : node.parent + '/' + node.text);
                            if (frm !== to) {
                                var form = $('<form method="post" class="hide" action="/project/' + project + '/file/move"><input name="to" value="' + to + '"><input name="from" value="' + frm + '"><input name="type" value="renam"></form>')
                                $('body').append(form);
                                form.submit();
                            }
                        });
                    }
                },
                "delete" : {
                    "separator_before"  : false,
                    "separator_after"   : false,
                    "label"             : "Delete",
                    "action"            : function(obj) {
                        var n = t.get_node(obj.reference[0].id);
                        var msg = 'Are you sure you would like to delete this file?';
                        if (n.type === 'dir') {
                            msg = 'Are you sure you would like to delete this folder and ALL of it\'s contents?';
                        }
                        swal({
                            title: '',
                            text: msg,
                            type: 'warning',
                            showCancelButton: true,
                            confirmButtonColor: 'red',
                            confirmButtonText: "Yes",
                            closeOnConfirm: false
                        }, function(){
                            var form = $('<form method="post" class="hide" action="/project/' + project + '/file/del"><input name="path" value="' + n.id + '"></form>')
                            $('body').append(form);
                            form.submit();
                        });
                    }
                }
            }
        },
        "plugins" : [
            "contextmenu", "dnd", "search",
            "state", "types", "wholerow", "sort"
        ]
    });
    t = $('#filetree').jstree();
});
