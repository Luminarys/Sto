$(document).ready(function()
{
    var theform = document.getElementById('mainform');
    var fileSelect = document.getElementById('files');
    var uploadButton = document.getElementById('sub');
    var fileDisplay = document.getElementById('fileDisplay');
    $('#button-about').attr('href', '#aboutmodal');
    $('#button-about').leanModal();

    if($('#files').prop('files').length > 0) {
        $('#sub').removeClass('disabled');
        $('#sub').addClass('waves-effect waves-light');
    } else {
        $('#sub').addClass('disabled');
        $('#sub').removeClass('waves-effect waves-light');
    }

    $('#files').click(function() {
        if(uploadInProgress) {
            return false;
        }
    });

    $('#files').change(function(){
        var fileSelect = document.getElementById('files');    
        if($('#files').prop('files').length > 0) {
                $('#sub').removeClass('disabled');
                $('#sub').addClass('waves-effect waves-light');
            } else {
                $('#sub').addClass('disabled');
                $('#sub').removeClass('waves-effect waves-light');
            }
        var delay = 0;
        if(onStandBy.length > 0) { delay = 500; }

        for (var i = 0; i < onStandBy.length; i++) {
            $('#file'+onStandBy[i]).fadeOut(500);
        };
        onStandBy.length = 0;
        setTimeout(function() {
        var _ltotal = total;
            for (var i = fileSelect.files.length - 1; i >= 0; i--) {
                var _id = i + _ltotal;
                fileAlertCard(fileSelect.files[i].name, 'on stand by', _id);
                $('#file' + _id + ' > div > div').removeClass('darken-1');
                $('#file' + _id + ' > div > div').addClass('lighten-2');
                onStandBy.push(_id);
            }
        }, delay);
    });

    theform.onsubmit = function(event) {
        event.preventDefault();
        var files = fileSelect.files;
        if(files.length > 0 && !uploadInProgress) {
            var ltotal = total;
            onStandBy.length = 0;
            $('#sub').addClass('disabled');
            $('#sub').removeClass('waves-effect waves-light');
            $('#select-button').removeClass('waves-effect waves-light');
            $('#select-button').addClass('disabled');
            uploadInProgress = true;
            for (i = files.length - 1; i >= 0; i--) {
                var _id = i + ltotal;
                var file = files[i];
                if(file.size > 50 * 1024 * 1024) {
                    $('#file' + _id + ' > div > div').removeClass('lighten-2');
                    $('#file' + _id + ' > div > div').addClass(ColorPalette[0] + ' darken-1');
                    $('#file' + _id + ' > div > div > .card-action > a').text('File too large!');
                } else {
                    $('#file' + _id + ' > div > div').removeClass('lighten-2');
                    $('#file' + _id + ' > div > div').addClass(ColorPalette[0] + ' darken-1');
                    $('#file' + _id + ' > div > div > .card-action > a').fadeOut(100);
                    $('#file' + _id + ' > div > div > .card-action > .progress').show();
                }
            }
            processFilesRecursively(files)
        }
    }
    document.getElementById('dropzone').addEventListener('dragenter', eventFileDropIn, false);
    document.getElementById('dropzone').addEventListener('dragleave', eventFileDropOut, false);
    document.getElementById('dropzone').addEventListener('dragover', eventFileDragEnter, false);
    document.getElementById('dropzone').addEventListener('drop', eventFileDropped, false);
});
//Card action -> #file > div > div > .card-action
var onStandBy = [];
var uploadInProgress = false;


var fadeTime = 2000;
var gCounter = 0;
var total = 0;

function processFilesRecursively(fileArray)
{
    if(gCounter >= fileArray.length) {
        uploadInProgress = false;
        $('#files').val('');
        $('#sub').addClass('disabled');
        $('#sub').removeClass('waves-effect waves-light');
        $('#select-button').addClass('waves-effect waves-light');
        $('#select-button').removeClass('disabled');
        gCounter = 0;
        total += 1;
        return 0;
    } else {
    var id = total;
    var file = fileArray[gCounter];
    if(file.size > 50 * 1024 * 1024)
    {
        gCounter += 1;
        total += 1;
        processFilesRecursively(fileArray);
        return 0;
    }

    var formData = new FormData();

    formData.append('files[]', file, file.name);

    $('#file' + id + ' > div > div > .card-action > .progress > div').removeClass('indeterminate');
    $('#file' + id + ' > div > div > .card-action > .progress > div').addClass('determinate');

    var request = $.ajax({
        url: '/api/upload',
        type: 'POST',
        async: true,
        data: formData,
        processData: false,
        contentType: false,
        xhr: function() {
            var myXhr = $.ajaxSettings.xhr();
            if(myXhr.upload){
                myXhr.upload.addEventListener('progress',function(e){
                    if(e.lengthComputable){
                        var max = e.total;
                        var current = e.loaded;
                        var Percentage = (current * 100)/max;
                        $('#file' + id + ' > div > div > .card-action > .progress > div').width(Percentage + '%');
                    }  
                }, false);
            }
            return myXhr;
        },
        complete: function(jqXHR, textStatus) {
            var resp = jqXHR.responseText;
            var code = JSON.parse(resp);
            var stat = code["success"];
            var file = code["files"][0];
            var feedback = '';
            if(stat) {
                feedback = '<span class="badge"><a target="_blank" href="/' + file["url"] +'">' + file["url"] + '</a></span><a>upload complete!</a>';
            }else{
                feedback = '<a>Invalid filename or banned extension.</a>';
            } 

            setTimeout(function() {
                var toinsert = $(feedback);
                toinsert.hide();
                toinsert.insertAfter('#file' + id + ' > div > div > .card-action > .progress');
                toinsert.fadeIn(1000);
            }, 990);
            $('#file' + id + ' > div > div > .card-action > .progress').fadeOut(1000);
            gCounter += 1;
            total += 1;
            processFilesRecursively(fileArray);
        },
        error: function(jqXHR, textStatus, errorThrown) {
            alertCard(errorThrown);
            gCounter += 1;
            total += 1;
            processFilesRecursively(fileArray);
        }
    });
    }
}
//Card action -> #file > div > div > .card-action
function fileAlertCard(filename, text, id) {
    var start = '<div class="row card-out" id="file' + id + '"><div class="col s12 l6 offset-l3 m8 offset-m2"><div class="card ';
    var color = ColorPalette[0] + ' ';
    var inter = 'darken-1 z-depth-2"><div class="card-content white-text">'
    var head = '<h5 class="truncate">';
    var name =  filename;
    var action = '</h5></div><div class="card-action">';
    var progress = '<div class="progress" style="display: none;"><div class="indeterminate"></div></div>'
    var progress2 = '<a>' + text + '</a></div>';
    var actionend = '</div>';
    var inter2 = '';
    var end = '</div></div></div>';

    var div = $(start +
                color +
                inter +
                head +
                name +
                action +
                progress +
                progress2 +
                actionend +
                inter2 +
                end
                );
    
    div.hide();
    div.insertAfter('#afterthis');
    div.fadeIn(fadeTime);
}

function alertCard(text) {
    var start = '<div class="row card-out"><div class="col s12 l6 offset-l3 m8 offset-m2"><div class="card ';
    var color = ColorPalette[0] + ' ';
    var inter = 'darken-1 z-depth-2"><div class="card-content white-text">'
    var head = '<h5 class="truncate">';
    var icon = '<i class="small mdi-alert-warning" style="margin-right: 2%;"></i>';
    var name =  text;
    var inter2 = '</h5></div>';
    var end = '</div></div></div>';

    var div = $(start +
                color +
                inter +
                head +
                icon +
                name +
                inter2 +
                end
                );
    
    div.hide();
    div.insertAfter('#afterthis');
    div.fadeIn(fadeTime);
    div.delay(4000).fadeOut(fadeTime);
}

function eventFileDropIn(e) {
    $('.fileupload > div').removeClass('darken-1');
    $('.fileupload > div').addClass('lighten-1');
    $('#dropzone > .card-title').text('Drop here!');
}

function eventFileDropOut(e) {
    $('.fileupload > div').removeClass('lighten-1');
    $('.fileupload > div').addClass('darken-1');
    $('#dropzone > .card-title').text('Fuwa.se');
}

function eventFileDropped(e) {
    e.stopPropagation();
    e.preventDefault();
    if(uploadInProgress) return false;
    $('.fileupload > div').removeClass('lighten-1');
    $('.fileupload > div').addClass('darken-1');
    $('#dropzone > .card-title').text('Fuwa.se');
    var files = e.dataTransfer.files;
    var ltotal = total;
    for (var i = 0; i < onStandBy.length; i++) {
        $('#file'+onStandBy[i]).remove();
    };
    onStandBy.length = 0;
    $('#sub').addClass('disabled');
    $('#sub').removeClass('waves-effect waves-light');
    $('#select-button').removeClass('waves-effect waves-light');
    $('#select-button').addClass('disabled');
    uploadInProgress = true;
    for (i = files.length - 1; i >= 0; i--) {
        var _id = i + ltotal;
        var file = files[i];
        fileAlertCard(files[i].name, 'on stand by', _id);
        if(file.size > 50 * 1024 * 1024) {
            $('#file' + _id + ' > div > div > .card-action > a').text('File too large!');
        } else {
            $('#file' + _id + ' > div > div > .card-action > a').fadeOut(100);
            $('#file' + _id + ' > div > div > .card-action > .progress').show();
        }
    }
    processFilesRecursively(files)
  }

  function eventFileDragEnter(e) {
    e.stopPropagation();
    e.preventDefault();
    $('.fileupload > div').removeClass('darken-1');
    $('.fileupload > div').addClass('lighten-1');
    $('#dropzone > .card-title').text('Drop here!');
  }
