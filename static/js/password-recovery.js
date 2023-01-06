$(document).ready(function(){
    M.AutoInit();
    Waves.displayEffect();

    $("#button").click(function(){
        if ($("#password").val() !== $("#confirmPassword").val()) {
            M.toast({html: "Passwords are different."});
            return;
        }

        M.toast({html: "Resetting password!"});
 
        $.ajax({
            url: "/password-recovery",
            type: "POST",
            contentType: "application/json; charset=utf-8",
            data: JSON.stringify({
                Code: getUrlParameter("code"),
                Password: $("#password").val(),
                Captcha: grecaptcha.getResponse()
            }),
            dataType: "json",
            success: function(r) {
                console.log(r.Code);
                switch(r.Code) {
                    case 0: {
                        window.location.replace("/panel");
                        break;
                    }
                    case 2: {
                        M.toast({html: "Please check the recaptcha and try again."});
                        break;
                    }
                    case 3: {
                 