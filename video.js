const video = document.getElementById("video");
const canvas = document.getElementById("canvas");
const snap = document.getElementById("snap");
const photo = document.getElementById("photo");
const context = canvas.getContext("2d");
const videoSelect = document.getElementById("videoSource");

// List cameras and microphones.
navigator.mediaDevices.enumerateDevices().then(gotDevices).catch(handleError);

// On change event for video source selection
videoSelect.onchange = getStream;

// Get available video sources
function gotDevices(deviceInfos) {
  for (let i = 0; i !== deviceInfos.length; ++i) {
    const deviceInfo = deviceInfos[i];
    const option = document.createElement("option");
    option.value = deviceInfo.deviceId;
    if (deviceInfo.kind === "videoinput") {
      option.text = deviceInfo.label || `camera ${videoSelect.length + 1}`;
      videoSelect.appendChild(option);
    }
  }
}

// Get stream based on selected video source
function getStream() {
  if (window.stream) {
    window.stream.getTracks().forEach((track) => {
      track.stop();
    });
  }
  const videoSource = videoSelect.value;
  const constraints = {
    video: { deviceId: videoSource ? { exact: videoSource } : undefined },
  };
  navigator.mediaDevices
    .getUserMedia(constraints)
    .then(gotStream)
    .catch(handleError);
}

function gotStream(stream) {
  window.stream = stream; // make stream available to console
  video.srcObject = stream;
}

function handleError(error) {
  console.error("Error: ", error);
}

// Capture the photo when the button is clicked
snap.addEventListener("click", () => {
  context.drawImage(video, 0, 0, canvas.width, canvas.height);
  const data = canvas.toDataURL("image/png");

  // Convert data URL to Blob
  const blob = dataURLToBlob(data);

  // Create form data
  const formData = new FormData();
  formData.append("file", blob, "photo.png");

  // Upload the image to the server
  fetch("http://localhost:8080/image/" + code, {
    method: "POST",
    body: formData,
  })
    .then((response) => response.json())
    .then((data) => {
      sec += 1
      console.log("Success:", data);
    })
    .catch((error) => {
      console.error("Error:", error);
    });
});

// Convert data URL to Blob
function dataURLToBlob(dataURL) {
  const parts = dataURL.split(";base64,");
  const byteString = atob(parts[1]);
  const mimeString = parts[0].split(":")[1];
  const ab = new ArrayBuffer(byteString.length);
  const ia = new Uint8Array(ab);
  for (let i = 0; i < byteString.length; i++) {
    ia[i] = byteString.charCodeAt(i);
  }
  return new Blob([ab], { type: mimeString });
}

// Get initial stream
getStream();
