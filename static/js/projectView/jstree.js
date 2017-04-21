$(document).ready(function() {
    $('#filetree').on('select_node.jstree', function(e, data) {
        /*console.log(e);
        console.log(data);
        e2 = e;
        evt = window.event || e;
        var button = evt.which || evt.button;
        // console.log(button);
        if ( button != 1 ) {
            return;
        }*/

        if (data.event != null) {
            if (data.event.type == 'contextmenu') {
                return
            }
        }

        var n = data.node;

        if (n.id == current && n.type != 'dir') {
            return
        }

        current = n.id;

        if (n.type === 'dir') {
            if (tree.is_closed(n)) {
                tree.open_node(n);
            } else {
                tree.close_node(n);
            }
            return;
        }

        $('div#fileName').text(n.text);

        getFile(n.id);

    }).on('move_node.jstree', function(e, data) {
        var frm = data.node.id;
        var to  = ((data.parent == "#") ? '/' + data.node.text : data.parent + '/' + data.node.text);
        if (frm !== to) {
            var form = $('<form method="post" class="hide" action="/project/' + project + '/file/move"><input name="to" value="' + to + '"><input name="from" value="' + frm + '"><input name="type" value="mov"></form>')
            $('body').append(form);
            form.submit();
        }
    }).jstree({
        "plugins" : [
            "contextmenu", "dnd", "search",
            "state", "types", "wholerow", "sort"
        ],
        "core" : {
            "multiple": false,
            "animation" : 0,
            "check_callback" : function(operation, node, node_parent, node_position, more) {
                if (operation === "move_node") {
                    return isSaved(node.id);
                }
                return true;
            },
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
                "valid_children" : ["dir", "file", "html", "css", "javascript"]
            },
            "dir" : {
                "icon" : "glyphicon glyphicon-folder-open",
                "valid_children" : ["dir", "file", "html", "css", "javascript"]
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
            "show_at_node":false,
            items : {
                "new" : {
                    "separator_before"  : false,
                    "separator_after"   : true,
                    "label"             : "New",
                    "action"            : false,
                    "submenu" :{
                        "create_file" : {
                            "label" : "File",
                            "icon": "fa fa-file-code-o",
                            action : function (obj) {
                                /*var n = tree.get_node(obj.reference[0].id);
                                while (n.type !== 'dir' && n.id !== '#') {
                                    n = tree.get_node(n.parent);
                                }
                                $('input#filePath').val(n.id);
                                $('div#newFileModal').modal('show');*/
                                newFile(obj.reference[0].id);
                            }
                        },
                        "create_folder" : {
                            "seperator_before" : false,
                            "seperator_after" : false,
                            "label" : "Folder",
                            "icon": "glyphicon glyphicon-folder-open",
                            action : function (obj) {
                                /*var n = tree.get_node(obj.reference[0].id);
                                while (n.type !== 'dir' && n.id !== '#') {
                                    n = tree.get_node(n.parent);
                                }
                                $('input#folderPath').val(n.id);
                                $('div#newFolderModal').modal('show');*/
                                newFolder(obj.reference[0].id);
                            }
                        }
                    }
                },
                "rename": {
                    "separator_before"  : false,
                    "separator_after"   : false,
                    "label"             : "Rename",
                    "action"            : function(obj) {
                        var n = tree.get_node(obj.reference[0].id);
                        if (!isSaved(n.id)) {
                            $.Notification.autoHideNotify('error', 'top center', 'You have unsaved changes.<br>Please save before renaming.');
                            return
                        }
                        var old_name = n.text;
                        tree.edit(n, 0, function(node, status, cancel) {
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
                        var n = tree.get_node(obj.reference[0].id);
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
                },
                "save" : {
                    "separator_before"  : true,
                    "separator_after"   : false,
                    "label"             : "Save",
                    "action"            : function(obj) {
                        var n = tree.get_node(obj.reference[0].id);
                        save(n.id);
                    }
                }
            }
        }
    });

    tree = $('#filetree').jstree();

    $(document).on('context_hide.vakata', function() {
        if (tree.get_selected()[0] != current) {
            tree.deselect_all();
            tree.select_node(current);
        }
    })
});
