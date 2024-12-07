let audioContext;
let analyser;
let microphone;
let javascriptNode;
let mediaRecorder;
let audioChunks = [];
let silenceTimer;
let isRecording = false;
const silenceThreshold = 7.5;
const silenceTimeout = 1000;

const startBtn = document.getElementById("startBtn");
const stopBtn = document.getElementById("stopBtn");
const status = document.getElementById("status");
const message = document.getElementById("msg");

startBtn.addEventListener("click", startDetection);
stopBtn.addEventListener("click", stopDetection);

function startDetection() {
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

        if (average > silenceThreshold) {
          console.log("Sound detected");
          status.textContent = "Status: Sound detected";
          silenceTimer = 0;
          if (!isRecording) {
            mediaRecorder.start();
            isRecording = true;
            console.log("Recording started");
          }
        } else {
          console.log("Silence detected");
          silenceTimer++;
          if (silenceTimer > silenceTimeout / 1000 && isRecording) {
            mediaRecorder.stop();
            isRecording = false;
            status.textContent = "Status: Silent";
            console.log("Recording stopped due to silence");
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
  startBtn.disabled = false;
  stopBtn.disabled = true;
  status.textContent = "Status: Inactive";
  clearTimeout(silenceTimer);
}

function sendAudioToServer(audioBlob) {
  const formData = new FormData();
  formData.append("file", audioBlob, "recording.wav");

  fetch("http://localhost:8080/upload", {
    method: "POST",
    body: formData,
  })
    .then((response) => response.json())
    .then((data) => {
      message.innerHTML += `<li>${data.message}</li>`;
    })
    .catch((error) => {
      console.error("Error sending audio to server", error);
    });
}
