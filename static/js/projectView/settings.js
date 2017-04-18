$(document).ready(function() {
    // configure setting on load

    $('#theme').val(editor.getTheme());
    $('#fontSize').val(editor.getFontSize());

    if (settings.editor.keys === 'vim') {
        $('#keybindings').val('vim');
    } else if (settings.editor.keys === 'emacs'){
        $('#keybindings').val('emacs');
    }
    if (settings.editor.load === 'load') {
        $('#load')[0].checked = true;
    }
    if (settings.version !== '' && settings.version !== undefined) {
        $('select#version').val(settings.version);
    }

    // update theme on change
    $('#theme').change(function() {
        var theme = $('#theme').val()
        editor.setTheme(theme);
        settings['editor']['theme'] = theme;
        saveSettings(settings);
    });

    // update font on change
    $('#fontSize').change(function() {
        var fontSize =+ $('#fontSize').val();
        editor.setFontSize(fontSize);
        settings['editor']['fontSize'] = fontSize;
        saveSettings(settings);
    });

    // update keybindings on change
    $('#keybindings').change(function() {
        var bind = $('#keybindings').val();
        if (bind === 'ace') {
            editor.setKeyboardHandler("");
        } else if (bind === 'vim') {
            editor.setKeyboardHandler("ace/keyboard/vim");
        } else {
            editor.setKeyboardHandler("ace/keyboard/emacs");
        }
        settings['editor']['keys'] = bind;
        saveSettings(settings);
    });

    // update 'Format on Run' on change
    $('#formatOnRun').change(function() {
        settings['formatOnRun'] = this.checked;
        saveSettings(settings);
    });

    // update load on change
    $('input[name="load"]').change(function() {
        settings['editor']['load'] = $('#load:checked').val();
        saveSettings(settings);
    });

    $('select#version').change(function() {
        settings['version'] = $('select#version').val();
        saveSettings(settings);
    });

    // reset settings to default on click (this includes favorites)
    $('#reset').click(function() {
        settings = {'editor':{}};
        saveSettings(settings);
        location.reload();
    });

});
