Detectword_pico: A Go Spoken Word Detector for the Raspi Pico - Ray Schuler 2022.04.14


Under Construction -- Come back soon.

Summary:

Detectword_pico is a system to compare spoken words with a pre-defined reference word, setting a logic output pin based on the detected word. When the system is powered on, the first two received words become the reference words. Subsequent words are compared to the reference words, and the output logic state is set accordingly. An example usage is learnig the words 'on' and 'off' to control a lamp, as shown in the following video.  The software is written in the Go (@go-ref) programming language and compiled with Tinygo (@go-ref).  The hardware target is a Raspberry Pi Pico (@pico-ref). The design attempts to acheive reasonable voice control, with minimal resources.

https://youtu.be/cquPffC5l68
Video demonstration: 'detectword_pico' controlling a lamp

Discussion
----------
The Detectword_pico project uses the Tinygo Go language compiler for embedded environments (@ref-tinygo) to create a word detector, commonly refered to as a 'hot word' or 'wake word' detector. The target hardware is the Raspberry Pi Pico, a low cost and high function ARM microcontroller with user programmable analog inputs and general purpose digital I/O.

Popular techniques for spoken word classification involve creating a spectrogram (@ref-spectrogram), which encapsulates time and frequency characteristics into a 2 dimensional scalar array.  These arrays are well suited to visual image representations, and image manipulation techniques.  The spectrogram in Figure (1) represents the word 'raspberry' generated from a 4096 voice sample captured by the Pico's analog to digital converter (ADC).  The colors represent the signal intensity in the indexed frequency bin (vertical), during the indexed time bin (horizontal). The colors range from blue to red, representing low and high intensity, respectively.  

![alt](https://github.com/schuler-robotics/detectword_pico/blob/master/images/xt-raspberry-4096-250.png)
![alt](https://github.com/schuler-robotics/detectword_pico/blob/master/images/spect-raspberry-4096-250-64.png)

Figure (1): A 4096 sample time domain wavform and spectrogram of the spoken word ‘raspberry’.

Image classification is often accomplished by machine learning techniques such as training a convolutional neural network (ConvNet) (ref-convnet). ConvNets are incredibly good at classifying images, including spectrograph images. A project affecting wake word detection on the Pico using machine learning techniques is listed in the references (@ref-ml-pico)  

The downside of training and using machine learning techniques for simple applications include high complexity, large processing and memory demands, and large sets of training data.  The machine learning model is typically trained on machine significantly more powerful than the microcontroller target. Detectword_pico does not rely on machine learning techniques, and instead uses simpler data reduction methods with memory and processing requirements suitable for low cost Pico microcontroller.

Popular 'smart speakers' rely on machine learning techques and send voice recordings to remote servers for processing, even when the device under control (e.g. a lamp) are inches from the speaker.  Sending audio data to servers outside of the end user's control involve privacy issues that may not be justified for single word voice control.

I set out to accomplish three objectives with Detectword_pico:

1) Determine if a reasonably accurate single word detector could be achieved soley on the Pico without the complexity and computational cost of machine learning techniques.

2) Learn the Go programming language, and assess the feasibility of using Go in the embedded application space.

3) Implement voice controlled lighting in my home, without relying on inernet based services.

The Process:
Upon power up, detectword_pico captures two reference words with one of the Pico's ADC.  These reference words are normalized in both amplitude and time, before being converted to spectrographs, and reduction techniques are applied.  The same process is applied to subsequent spoken words, and a square error of the reduced data between the reference and target word is calculated.  This squared error scalar value is used to predict if a target word matches one of the refernce words.  The 0V-3.3V logic state of a GPIO pin tracks last detected reference word, for example 'on' or 'off'.  When neither reference word is detected, the GPIO state remains unchanged.

The detectword.go function CreateU16SpectFromU16() converts the voice samples into a 2 dimensional spectrograph array. The input waveform of 'buf_size' samples is normalized and broken into 'Tbins' time segments. Each time segment is filtered by a Hamming window (@ref-hamming) to suppress the discontinuities created by segmenting the data. The filtered segments are then converted to the frequency domain, generating 'Fbins' values for the spectgrograph. It was important to use an 'in place' discrete fourier transform (DFT) algorithm (@ref-tukey) to conserve memory on the Pico.  Detectword pico relys on the FOSS go-fft package (@ref-go-fft)to generate DFTs.

The data reduction consists of a two step pooling process to which reduces the memory and processing requirements for comparing words.  The first pooling operation steps a window with size determined by 'block' parameters, and creates a reduced two dimensional array with elements equal to the average value of each window position.  The second stage repeats the pooling process returning the peak value of a smaller window stepped across the result of average pooling.

The default tuning parameters in Detectword_pico include a frequency bin range of approximately -1.8Khz to 1.8Khz, with time bins ranging from 0sec to approximately 300msec.  The block sizes for average and peak pooling are 8x8 and 4x4, respectively. These value are adequate for single sylable human voice word detection.

Figures (2-4) show spectrograms, and their reduced counterparts, for the words 'on' and 'off'.  The final reduced arrays in Figure (4) are used to determine if a target word closely resembles a reference word. The images of the spectrograms and reduced arrays encode increasing frequency power values ranging from blue to red. 

<img src="https://github.com/schuler-robotics/detectword_pico/blob/master/images/spect-on-4096-250-64.png" width="400" height="300">

![alt](https://github.com/schuler-robotics/detectword_pico/blob/master/images/spect-on-4096-250-64.png)
![alt](https://github.com/schuler-robotics/detectword_pico/blob/master/images/spect-off-4096-250-64.png)

Figure (2) Spectrographs of the words 'on' and 'off'; 4096 samples, 64 time bins, and 64 frequency bins

![alt](https://github.com/schuler-robotics/detectword_pico/blob/master/images/pool-avg-on-4096-250-64.png)
![alt](https://github.com/schuler-robotics/detectword_pico/blob/master/images/pool-avg-off-4096-250-64.png)

Figure (3) The average pooling spectrographs of 'on' and 'off', reduced from Figure (1)
![alt](https://github.com/schuler-robotics/detectword_pico/blob/master/images/pool-peak-on-4096-250-64.png)
![alt](https://github.com/schuler-robotics/detectword_pico/blob/master/images/pool-peak-off-4096-250-64.png)

Figure (4) The peak value pooling spectrographs 'on' and 'off', reduced from Figure (3)

While the Pico has sufficient memory to set Detectword_pico's 'buf_size' to 4096 samples, the response time is large-- on the order of 2 seconds.  As part of the tuning process, I determined that 1024 samples provide reasonaby good detection with a greatly reduced response time.  The response time to process single sylable words from 1024 ADC samples is in the hundreds of milliseconds range. The demonstration video above is processing 1024 samples voice waveforms.

Voice acquisition for the Pico is obtained with a peizo microphone and a Maxim 4466 amplifier feeding the Pico ADC. In addition to the tuning parameters previoiusly described, Detectword_pico includes threshold parameters to gate capture.  Storing voice samples begins when the ADC detects a sound level above the software parameter 'threshold'. The capture continues until 'buf_size' samples have been collected. The end of the capture buffer is also truncated of sounds below 'threshold'.  Removing "quiet" samples from the end of the buffer allows the time domain waveform to be normalized in both amplitude and time.  Another threshold paramater, 'SpectThresh', limits

The paramters used to tune word detection are: capture sample size in bytes (buf_size), ADC sampling rate (Tsamp), number of spectrogram time bins (Tbins), number of spectrogram frequency bins (Fbins), pooling block sizes, and noise thresholds (threshold and SpectThresh).

Possible Improvements
---------------------

Word detection would likely be improved by low pass filtering the microphone ADC input to limit frequencies above the Nyquest rate.  Empirically, I found most meaningful data from my voice between 200Hz-400Hz. The generated spectrograms currently include some aliasing, this folded frequency content has not prevented reasonable word detection.

Adding automatic gain control the microphone would improve detection performance for words spoken at different distances or loudness than the recorded reference words.

Generating a spectrograph is a parallel process, and the Pico has two cores.  Breaking the spectrograph genration into two concurrent processes will reduce the response time.

Negative spectrogram frequencies are maintained in detectword_pico for image aesthetics.  These frequencies are not included in pooling reduction.  Removing negative frequencies from the spectrograms will improve processing speed and reduce memory usage.

Allowing the spectrogram time bins to overlap in would reintroduce valid detection data suppressed by the Hamming filter, at the expense of increased response time.

Conclusions:

While ConvNets and other machine learning techniques provide powerful tools speech detection, this project shows they are not strictly required for simple word detection on a low cost microcontroller.  The techniques employed by Detectword_pico provide a reasonably good solution for a voice controlled lamp.

The Go language (ref-go), and the memory efficient Tinygo compiler (ref-tinygo), exceeded my expectations for developing on the Pico.  The code is written with standard Go libraries, with the exception of the memory efficient "in place" DFT implementaiton from the go-fft package (ref-go-fft).

Software Details:

Tinygo v0.21 which is based on Go v1.17.6
'detectword_pico.go' is the main() entry point of the program
'detectword.go' includes functions specific to the detectword application.  
'utils_dw.go' includes functions applicable to wider range of DSP applications. 
'fft.go' and 'errors.go' are manually included from the go-fft package, as Tinygo v0.21 does not support all dependencies.

--
(@ref-go) https://go.dev
(@ref-tinygo) https://github.com/tinygo-org
(@ref-spectrogram) https://en.wikipedia.org/wiki/Spectrogram
(@ref-convnet) https://developers.google.com/machine-learning/practica/image-classification/convolutional-neural-networks
(@ref-ml-pico) https://github.com/henriwoodcock/pico-wake-word
(@ref-hamming) https://stackoverflow.com/questions/5418951/what-is-the-hamming-window-for
(@ref-tukey) https://en.wikipedia.org/wiki/Cooley%E2%80%93Tukey_FFT_algorithm
(@ref-go-fft) https://github.com/ledyba/go-fft/blob/master/LICENSE


