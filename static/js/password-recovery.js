$(document).ready(function(){
    M.AutoInit();
    Waves.displayEffect();

    $("#button").click(function(){
        if ($("#password").val() !== $("#confirmPassword").val()) {
            M.toast({html: "Passwords are different."});
            return;
        }

        M.toast({html: "Resettin