import { useState } from "react";
import axios from "axios";

import "./App.css";

interface FileInputProps {
  file: File | null;
  setFile: React.Dispatch<React.SetStateAction<File | null>>;
}

interface PresignedUrlResponse {
  key: string;
  get: string;
  put: string;
}

interface UploadRequest {
  Key: string;
}

function App() {
  const [file, setFile] = useState<File | null>(null);
  const [success, setSuccess] = useState(false);
  const [accessUrls, setAccessUrls] = useState<string[]>([]);

  async function getPresignedUrls(files: UploadRequest[]) {
    //if files is empty, return
    if (files.length === 0) {
      return [];
    }
    try {
      const response = await axios.post<PresignedUrlResponse[]>(
        "http://localhost:1323/request-presigned-url",
        files
      );
      return response.data;
    } catch (error) {
      console.error(error);
      throw error;
    }
  }

  const handleChange = (e: any) => {
    setFile(e.target.files[0]);
  };

  const handleUpload = async () => {
    // Example code to handle file upload
    setSuccess(false);
    const uploadRequests = file ? [{ Key: file.name }] : []; // this is not optinal since the key is file name. This should be
    const presignedUrls = await getPresignedUrls(uploadRequests);
    //if the presignedUrls is empty, return
    if (presignedUrls.length === 0) {
      setSuccess(false);
      return;
    }

    //with each presigned url, upload the file
    const uploadPromises = presignedUrls.map((presignedUrls) => {
      return axios.put(presignedUrls.put, file, {
        headers: {
          "Content-Type": file?.type,
        },
      });
    });

    //wait for all uploads to complete
    await Promise.all(uploadPromises)
      .then((res) => console.log(res))
      .catch((error) => {
        console.error(error);
        setSuccess(false);
        throw error;
      });
    // return access urls after all uploads are complete
    const accessUrls = presignedUrls.map((presignedUrls) => {
      return presignedUrls.get;
    });
    setSuccess(true);
    setAccessUrls(accessUrls);
    console.log(accessUrls);
  };

  return (
    <div>
      <h1>Upload File</h1>
      {success && <p>File uploaded successfully</p>}
      <input type="file" onChange={handleChange} />
      <button onClick={handleUpload}>Upload</button>
      {file && <p>Selected file: {file.name}</p>}
      {accessUrls.length > 0 && (
        <div>
          <h2>Access URLs</h2>

          {accessUrls.map((url) => (
            <>
              <h3>{url}</h3>
              <img
                src={url}
                alt="uploaded"
                style={{ width: "300px", height: "300px" }}
              />
            </>
          ))}
        </div>
      )}
    </div>
  );
}

export default App;
