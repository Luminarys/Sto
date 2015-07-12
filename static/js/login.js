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
                var resp = JSON.parse(data);
                if(resp.success){
                    alert("Success!");
                }else{
                    $("#error-msg").remove();
                    $(".card-content").prepend("<span id='error-msg' class='red-text text-darken-1'>" + resp.message + "</span>");

                }
            }
        });
    });
});
