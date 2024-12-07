let audioContext;
let analyser;
let microphone;
let javascriptNode;
let mediaRecorder;
let audioChunks = [];
let silenceTimer;
let isRecording = false;
const silenceThreshold = 7.0;
const silenceTimeout = 200;

const startBtn = document.getElementById("startBtn");
const stopBtn = document.getElementById("stopBtn");
const status = document.getElementById("status");
const message = document.getElementById("msg");

const videoQueue = [];
let currentVideoIndex = 0;
const videoPlayer = document.getElementById("video-player");

startBtn.addEventListener("click", startDetection);
stopBtn.addEventListener("click", stopDetection);

let sec = 0;
let code = "";

function loadNextVideo() {
  if (currentVideoIndex < videoQueue.length) {
    videoPlayer.src = "kamus/" + videoQueue[currentVideoIndex] + ".webm";
    videoPlayer.playbackRate = 3.0;
    videoPlayer.play();
  }
}

videoPlayer.addEventListener("ended", () => {
  currentVideoIndex++;
  if (currentVideoIndex < videoQueue.length) {
    loadNextVideo();
  } else {
    console.log("Semua video telah selesai diputar.");
  }
});

function startDetection() {
  generateCode();
  navigator.mediaDevices
    .getUserMedia({ audio: true })
    .then((stream) => {
      audioContext = new (window.AudioContext || window.webkitAudioContext)();
      analyser = audioContext.createAnalyser();
      microphone = audioContext.createMediaStreamSource(stream);
      javascriptNode = audioContext.createScriptProcessor(2048, 1, 1);

      analyser.smoothingTimeConstant = 0.8;
      analyser.fftSize = 1024;

      mediaRecorder = new MediaRecorder(stream);
      mediaRecorder.ondataavailable = (event) => {
        audioChunks.push(event.data);
      };

      mediaRecorder.onstop = () => {
        const audioBlob = new Blob(audioChunks, { type: "audio/wav" });
        audioChunks = [];
        sendAudioToServer(audioBlob);
      };

      microphone.connect(analyser);
      analyser.connect(javascriptNode);
      javascriptNode.connect(audioContext.destination);

      javascriptNode.onaudioprocess = function () {
        const array = new Uint8Array(analyser.frequencyBinCount);
        analyser.getByteFrequencyData(array);
        const arraySum = array.reduce((a, b) => a + b, 0);
        const average = arraySum / array.length;
        console.log(average, "<><><>", silenceThreshold);
        if (average > silenceThreshold) {
          status.textContent = "Status: Sound detected";
          silenceTimer = 0;
          if (!isRecording) {
            mediaRecorder.start();
            isRecording = true;
          }
        } else {
          silenceTimer++;
          if (silenceTimer > silenceTimeout / 1000 && isRecording) {
            mediaRecorder.stop();
            isRecording = false;
            status.textContent = "Status: Silent";
          }
        }
      };

      startBtn.disabled = true;
      stopBtn.disabled = false;
      status.textContent = "Status: Detecting...";
    })
    .catch((error) => console.error("Error accessing media devices.", error));
}

function stopDetection() {
  if (audioContext) {
    audioContext.close();
    audioContext = null;
  }
  if (isRecording) {
    mediaRecorder.stop();
    isRecording = false;
  }
  getSummary()
  startBtn.disabled = false;
  stopBtn.disabled = true;
  status.textContent = "Status: Inactive";
  clearTimeout(silenceTimer);
}

function sendAudioToServer(audioBlob) {
  const formData = new FormData();
  formData.append("file", audioBlob, "recording.wav");

  fetch(`http://backendwppdev.my.id/audio/${sec}/${code}`, {
    method: "POST",
    body: formData,
  })
    .then((response) => response.json())
    .then((data) => {
      message.innerHTML += `<li class="bg-gray-100 p-4 rounded-md shadow" >${data.message}</li>`;
      videoQueue.push(...data.video_queue);
      console.log(data.video_queue)
      // videoQueue.push(...data.message?.toLowerCase().split(" "));
      loadNextVideo();
    })
    .catch((error) => {
      console.error("Error sending audio to server", error);
    });
}

function generateCode() {
  fetch("http://backendwppdev.my.id/start", {
    method: "GET",
  })
    .then((response) => response.json())
    .then((data) => {
      sec = data.sec;
      code = data.code;
      document.getElementById("code").innerHTML = data.code;
    })
    .catch((error) => {
      document.getElementById("container").innerHTML = `<h1>Internal Server Error</h1>`
      console.error("Error generate code", error);
    });
}

function getSummary() {
  fetch("http://backendwppdev.my.id/summary/" + code, {
    method: "GET",
  })
    .then((response) => response.json())
    .then((data) => {
      document.getElementById("code").innerHTML = data.code;
    })
    .catch((error) => {
      document.getElementById("container").innerHTML = `<h1>Internal Server Error</h1>`
      console.error("Error generate code", error);
    });
}