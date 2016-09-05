function onSignIn(googleUser) {
  var profile = googleUser.getBasicProfile();
  console.log('ID: ' + profile.getId()); // Do not send to your backend! Use an ID token instead.
  console.log('Name: ' + profile.getName());
  console.log('Image URL: ' + profile.getImageUrl());
  console.log('Email: ' + profile.getEmail());

  document.cookie = "id="+ profile.getId();
}

var upload =
{
    pageName: '',
    bytesPerChunk: 20 * 1024 * 1024,
    loaded: 0,
    total: 0,
    file: null,
    fileName: "",

    uploadFile: function () {
        var size = upload.file.size;

        if (upload.loaded > size) return;

        var end = upload.loaded + upload.bytesPerChunk;
        if (end > size) { end = size; }

        var blob = upload.file.slice(upload.loaded, end);

        var xhr = new XMLHttpRequest();

        xhr.open('POST', upload.pageName, false);

        xhr.setRequestHeader("Content-Type", "multipart/form-data");
        xhr.setRequestHeader("X-File-Name", upload.fileName);
        xhr.setRequestHeader("X-File-Type", upload.file.type);

        xhr.send(blob);

        upload.loaded += upload.bytesPerChunk;

        setTimeout(upload.updateProgress, 100);
        setTimeout(upload.uploadFile, 100);
    },
    upload: function (file) {
        upload.file = file;

        var date = new Date();
        upload.fileName = date.format("dd.MM.yyyy_HH.mm.ss") + "_" + file.name;

        upload.loaded = 0;
        upload.total = file.size;

        setTimeout(upload.uploadFile, 100);


        return upload.fileName;
    },
    updateProgress: function () {
        var progress = Math.ceil(((upload.loaded) / upload.total) * 100);
        if (progress > 100) progress = 100;

        $("#dvProgressPrcent").html(progress + "%");
        $get('dvProgress').style.width = progress + '%';
    }
};