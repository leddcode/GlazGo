> ⚠️ **Disclaimer**
>
> The development of **GlazGo** has progressed significantly beyond the version previously uploaded to this GitHub repository.
>
> At this time, the **latest source code will not be shared publicly**. Only compiled `.exe` releases will be made available for download through the official release section.
>
> This approach allows for cleaner releases while development continues in a more experimental and evolving form.
>
> Thanks for checking it out and supporting this early-stage project.


# GlazGo
Introducing a powerful fuzzing tool for pentesters and security researchers to investigate web applications and APIs. This tool is a compiled executable file that can be used when other tools are unavailable, such as when testing a web application in a black-box environment on a corporate machine. It helps to invetigate web apps and APIs by automating the process of providing expected and unexpected inputs in order to uncover application resources or to cause any unexpected behavior or crashes. 

# Usage
To use GlazGo, you will need to provide the URL of the website or application that you want to test.
The URL should contain the string "FUZZ" where you want the tool to inject test data.

You can also provide additional headers and cookies to include in the request.

To start the fuzzing process, you will need to run the tool and provide it with a list of test data to use.
This can be a list of words, numbers, or special characters. The tool will then send requests to the URL with
the test data injected at the "FUZZ" location and analyze the response for any errors or vulnerabilities.

For example, if you have a URL of "https://example.com/dir/FUZZ" and a list of test data that includes the words "test1" and "test2"
the tool will send requests to "https://example.com/dir/test1" and "https://example.com/dir/test2" and look for any issues in the response.

# Disclaimer
Using GlazGo for fuzzing targets without obtaining prior consent is illegal. The user assumes all responsibility for their use of the tool.

Keep in mind that fuzzing can generate a large number of requests and may potentially cause issues with the website or application being tested.
It is important to use caution and obtain permission before fuzzing any production systems.
