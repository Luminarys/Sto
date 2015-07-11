$(document).ready(function() {
    $("form#login").submit(function(e) {
        e.preventDefault();
        var formData = new FormData();
        var user = $("#uname").val();
        var password = $("#pwd").val();
        formData.append("username", user);
        formData.append("password", password);
        $.ajax({
            url: "login",
            type: "POST",
            data: formData,
            async: true,
            processData: false,
            contentType: false,
            success: function(data) {
                alert(data);
            }
        });
    });
});
