Detectword_pico: A Go Spoken Word Detector for the Raspi Pico
-------------------------------------------------------------
2022.04.18
<br>
<pre>
Reviewer notes: (@ref-tags) will become reference numbers (1), (2), ...
</pre>
The Detectword_pico application compares spoken words with predefined reference words, and sets a logic output pin based on the detected word. When the system is powered on, the first two received words become the reference words. Subsequent words are compared to the reference words controlling the output accordingly. The breadboard image below links a demonstration video using words 'on' and 'off' to control a lamp.  Detectword_pico is written in the Go (@go-ref) programming language and compiled with Tinygo (@ref-tinygo).  The hardware target is a Raspberry Pi Pico RP2040 (@ref-pico). The design attempts to achieve reasonable voice control, with minimal resources.

<p align="center">
<a href="https://youtu.be/cquPffC5l68" title="Video demonstration"><img src="https://img.youtube.com/vi/cquPffC5l68/maxresdefault.jpg" width="350px"/></a>
</p>


Discussion
----------
Detectword_pico is a word detector, also referred to as a 'hot word' or 'wake word' detector. The hardware target is the Raspberry Pi Pico board.  The Pico RP2040 is a high function ARM microcontroller with analog inputs, general purpose digital I/O, and a retail cost of 1 USD (2022). Detectword_pico is written in Go (@ref-go) and compiled to a UF2 firmware file with Tinygo (@ref-tinygo), a Go compiler for embedded environments. 

Spectrograms (@ref-spectrogram) are created encapsulating time and frequency features into two dimensional arrays.  These arrays are well suited to image processing techniques.  The spectrogram in Figure (1) represents the word 'raspberry' as captured by the Pico analog to digital converter (ADC).  Spectrogram colors represent frequency amplitudes (vertical), and duration (horizontal). Ranging from blue to red the colors represent low and high intensity, respectively.  

<p float="left">
<img src="https://github.com/schuler-robotics/detectword_pico/blob/master/images/xt-raspberry-4096-250.png" width="400" height="300" />
<img src="https://github.com/schuler-robotics/detectword_pico/blob/master/images/spect-raspberry-4096-250-64.png" width="400" height="300" />
</p>
Figure (1): A 4096 sample time domain waveform and spectrogram of the spoken word ‘raspberry’.
<br />
<br />

The use of spectrograms turns word detection into an image classification problem, often solved with machine learning (ML) techniques like convolutional neural networks (ConvNet) (ref-convnet). ConvNets are incredibly good at classifying images. A wake word detection project on the Pico using machine learning techniques is listed in the references (@ref-ml-pico).

The drawbacks of ML solutions include complexity, processing requirements for training operation, and training data requirements. ML models are trained on machines significantly more powerful than the Pico, and transferred to the target hardware.

Detectword_pico uses data reduction and comparison to affect simple word detection with less complexity and lower resource requirements than ML solutions.  

Popular 'smart speakers' rely on machine learning techniques and process voice recordings on remote servers, even when the device under control (e.g. a lamp) is inches away.  Sending voice data to servers outside of the end user's control involves privacy issues that may not be justified for single word voice control.

I set out to accomplish three objectives with Detectword_pico:

1) Determine if a reasonably accurate single word detector could be achieved on the Pico, without the complexity, computational cost, and external resources required for ML techniques.

2) Learn the Go programming language, and assess the feasibility of using Go in the embedded application space.

3) Implement low cost voice controlled lighting in my home, without relying on internet based services.

The Process
-----------
Upon powering up, detectword_pico captures two reference words with one of the Pico's ADC.  These reference words are normalized in both amplitude and time, before being converted to spectrograms and reduction techniques are applied.  The same process is applied to subsequent spoken words, and a sum squared error of the reduced data between the reference and target word is calculated.  This sum squared error is used to predict if the target word matches one of the reference words.  The 0V-3.3V logic state of a GPIO pin tracks the last detected reference word. For example 3.3V for 'on', and 0V for 'off'.  When neither reference word is detected, the GPIO state remains unchanged.

The detectword.go function CreateU16SpectFromU16() converts voice samples into a two  dimensional spectrogram array. An input waveform of 'buf_size' samples is normalized and broken into 'Tbins' time segments. Each time segment is filtered by a Hamming window (@ref-hamming) to suppress the discontinuities created by segmentation. The filtered segments are converted to the frequency domain, generating 'Fbins' values for the spectrogram. It was important to use an 'in place' discrete fourier transform (DFT) algorithm (@ref-tukey) to conserve memory on the Pico.  'In place' calculation means the time samples are presented to the DFT algorithm as real floating point values in a complex128 array, and are swapped out with frequency domain results in that same allocated memory. Detectword pico relies on the FFT() function from the very capable FOSS go-fft package (@ref-go-fft).

Figure (2) shows the calculated spectrograms for the words 'on' and 'off'.  Image representations of spectrograms and reduced arrays in this write up encode increasing frequency amplitudes as colors ranging from blue to red. 

<p float="left">
<img src="https://github.com/schuler-robotics/detectword_pico/blob/master/images/spect-on-4096-250-64.png" width="400" height="300" />
<img src="https://github.com/schuler-robotics/detectword_pico/blob/master/images/spect-off-4096-250-64.png" width="400" height="300" />
</p>
Figure (2) Spectrograms of the words 'on' and 'off'; 4096 samples, 64 time bins, and 64 frequency bins
<br />
<br />

Data reduction of the spectrograms consists of a two stage pooling process, reducing the memory and processing requirements of comparison.  The first pooling operation steps a rectangular window with a size determined by 'block' parameters across the spectrogram, creating a reduced two dimensional array with elements equal to the average value of spectrogram elements under that window position.  The second pooling stage repeats the process returning the peak value of a smaller window stepped across the array which resulted from average pooling. Pooling differs from convolution in that the window positions do not overlap.

Figure (3) shows the reduced spectrograms for the 'on' and 'off' spectrograms of Figure (2).  The smaller peak value arrays are used to predict how closely the target word matches a reference word.

<p float="left">
<img src="https://github.com/schuler-robotics/detectword_pico/blob/master/images/pool-avg-on-4096-250-64.png" width="200" height="200" />
<img src="https://github.com/schuler-robotics/detectword_pico/blob/master/images/pool-peak-on-4096-250-64.png" width="150" height="150" />
<img src="https://github.com/schuler-robotics/detectword_pico/blob/master/images/pool-avg-off-4096-250-64.png" width="200" height="200" />
<img src="https://github.com/schuler-robotics/detectword_pico/blob/master/images/pool-peak-off-4096-250-64.png" width="150" height="150" />
</p>
Figure (3) Average and peak pooling results from the spectrograms representing words 'on' and 'off'.
<br />
<br />

The default tuning parameters in Detectword_pico include a frequency bin range of approximately -1.8Khz to 1.8Khz, with time bins ranging from 0sec to approximately 300msec.  The block sizes for average and peak pooling are 8x8 and 4x4, respectively. These values are adequate for single syllable human voice word detection.

While the Pico has sufficient memory to set Detectword_pico's 'buf_size' to 4096 samples, the response time is large-- on the order of 2 seconds.  Setting 'buf_size' to 1024 samples provides reasonably good detection with a greatly reduced response time.  The response time to process single syllable words from 1024 ADC samples is in the hundreds of milliseconds range. The demonstration video above processes a 1024 sample voice waveform.  Figure (4) shows spectrograms from 1024 sample 'on' and 'off' voice captures.

<p float="left">
<img src="https://github.com/schuler-robotics/detectword_pico/blob/master/images/spect-on-1024-250-16.png" width="400" height="300" />
<img src="https://github.com/schuler-robotics/detectword_pico/blob/master/images/spect-off-1024-250-16.png" width="400" height="300" />
</p>
Figure (4) Spectrograms of the words 'on' and 'off'; 1024 samples, 16 time bins, and 16 frequency bins
<br />
<br />

Voice waveforms for Detectword_pico are obtained with a piezo microphone feeding a Maxim 4466 amplifier, with output tied to the Pico ADC. In addition to the tuning parameters previously described, Detectword_pico includes threshold parameters to gate capture and remove leading and trailing "quiet" periods.  Capture begins when the ADC detects a sound level above the software parameter 'threshold'. The capture continues until 'buf_size' samples have been collected. The end of the capture buffer is truncated of sounds below 'threshold'.  Removing "quiet" samples from the beginning and end of the buffer allows the sample waveform to be normalized in both amplitude and time, improving reference to target comparisons.  Another threshold parameter, 'SpectThresh', limits low amplitude noise in the spectrograms.

The parameters used to tune word detection are capture sample size in bytes (buf_size), ADC sampling rate (Tsamp), number of spectrogram time bins (Tbins), number of spectrogram frequency bins (Fbins), pooling block sizes, and noise thresholds (threshold and SpectThresh).

The complete Detectword_pico process flow, at greatly exaggerated scale, is illustrated in Figure (5). Each reference and target word undergoes the process, and the decision is based on the sum squared error of the final peak pooling stages.

<p float="left">
<img src="https://github.com/schuler-robotics/detectword_pico/blob/master/images/dw-process-sketch-20220419.png" width="800" />
</p>
Figure (5) Detectword_pico process diagram
<br />
<br />


Possible Improvements
---------------------
A single pole low pass filter (LPF) between amplifier and ADC suppresses noise and provides dc isolation and level shifting. Some frequency aliasing exists above the Nyquest rate. A higher order LPF would likely improve word detection. Empirically, I found most meaningful data from my voice is between 200Hz-400Hz. The level of aliasing, visible in the spectrograms above, has not prevented reasonable word detection.

Adding automatic gain control to the amplifier would improve detection performance of words spoken at different distances or loudness than the recorded reference words.

Generating a spectrogram is a parallel process, and the Pico has two cores.  Breaking spectrogram construction into two concurrent processes will reduce the response time.

Negative spectrogram frequencies are maintained in memory for image aesthetics only. Negative frequencies are not included in reduction.  Removing these frequencies from spectrogram generation would reduce response time and memory usage.

Allowing spectrogram time bins to overlap would reintroduce valid detection data suppressed by the Hamming filter. The overlaps would improve the spectrogram fidelity, at the expense of increased memory use and processing time.

Conclusions
-----------
This project attempts to show that simpler speech detection techniques than machine learning exist, specifically for word detection on low cost microcontrollers.  The techniques employed by Detectword_pico provide a reasonably good solution for a voice controlled lamp.

The Go language (ref-go), and Tinygo compiler (ref-tinygo) are capable and easy to learn tools for embedded systems development. Detectword_pico is written with standard Go libraries, with the exception of the fast and efficient DFT implementation from the go-fft package (ref-go-fft).

Logistics
---------
The Tinygo v0.21 compiler is based on Go v1.17.6.  'detectword_pico.go' is the main() entry point of the program.  'detectword.go' includes functions specific to the detectword application.  'utils_dw.go' includes functions applicable to a wider range of DSP applications.  'fft.go' and 'errors.go' are manually included from the go-fft package, as Tinygo v0.21 does not support all dependencies.

<pre>
Ray Schuler
schuler at usa.com
</pre>


References
----------
<p>
(@ref-go) https://go.dev<br>
(@ref-tinygo) https://github.com/tinygo-org<br>
(@ref-pico) https://www.raspberrypi.com/products/rp2040/
(@ref-spectrogram) https://en.wikipedia.org/wiki/Spectrogram<br>
(@ref-convnet) https://developers.google.com/machine-learning/practica/image-classification/convolutional-neural-networks<br>
(@ref-ml-pico) https://github.com/henriwoodcock/pico-wake-word<br>
(@ref-hamming) https://stackoverflow.com/questions/5418951/what-is-the-hamming-window-for<br>
(@ref-tukey) https://en.wikipedia.org/wiki/Cooley%E2%80%93Tukey_FFT_algorithm<br>
(@ref-go-fft) https://github.com/ledyba/go-fft/blob/master/LICENSE<br>
</p>
