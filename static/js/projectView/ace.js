
$(document).ready(function() {
    // setup ace editor
    editor.session.setMode("ace/mode/html");
    editor.setTheme('ace/theme/monokai');
    editor.renderer.setShowGutter(true);
    editor.setHighlightActiveLine(true);
    editor.setReadOnly(false);
    editor.setDisplayIndentGuides(true);
    editor.setOption("scrollPastEnd", 1);
    editor.renderer.setScrollMargin(15)
    editor.$blockScrolling = Infinity
    $('textarea.ace_text-input').focus();

    // init settings from local storage
    if (settings.editor !== undefined && !$.isEmptyObject(settings.editor)) {
        editor.setTheme(((settings.editor.theme === '' || settings.editor.theme === undefined) ? 'ace/theme/monokai' : settings.editor.theme));
        editor.setFontSize(((settings.editor.fontSize === '' || settings.editor.fontSize === undefined) ? 12 : settings.editor.fontSize));
        if (settings.editor.keys === 'vim') {
            editor.setKeyboardHandler("ace/keyboard/vim");
        } else if (settings.editor.keys === 'emacs'){
            editor.setKeyboardHandler("ace/keyboard/emacs");
        }
    } else {
        editor.setTheme('ace/theme/monokai');
        editor.setFontSize(12);
    }

    // add key bindings to ace editor
    editor.commands.addCommand({
        name: 'moveLineDown',
        bindKey: {win: 'Ctrl-down',  mac: 'Command-down'},
        exec: function(editor) {
            editor.moveLinesDown();
        },
        readOnly: false
    });

    editor.commands.addCommand({
        name: 'moveLineUp',
        bindKey: {win: 'Ctrl-up',  mac: 'Command-up'},
        exec: function(editor) {
            editor.moveLinesUp();
        },
        readOnly: false
    });

    editor.commands.addCommand({
        name: 'moveCursorUp',
        bindKey: {win: 'Alt-Shift-up',  mac: 'Alt-Shift-up'},
        exec: function(editor) {
            editor.selectMoreLines(-1);
        },
        readOnly: false
    });

    editor.commands.addCommand({
        name: 'moveCursorDown',
        bindKey: {win: 'Alt-Shift-down',  mac: 'Alt-Shift-down'},
        exec: function(editor) {
            editor.selectMoreLines(1);
        },
        readOnly: false
    });

    editor.commands.addCommand({
        name: "replace",
        bindKey: {win: "Ctrl-R", mac: "Command-Option-R"},
        exec: function(editor) {
            ace.config.loadModule("ace/ext/searchbox", function(e) {e.Search(editor, true)});
        }
    });
});
