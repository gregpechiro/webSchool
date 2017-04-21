$(document).ready(function() {
    $('.away').click(function() {
        var parent = $(this).closest('.side');
        if (parent.hasClass('menu-open')) {
            parent.find('.content').addClass('hide');
            parent.find('.alt').removeClass('hide');
            parent.removeClass('menu-open');
            parent.addClass('menu-close');
            return
        }
        parent.find('.alt').addClass('hide');
        parent.find('.content').removeClass('hide');
        parent.removeClass('menu-close');
        parent.addClass('menu-open');
    });
});
