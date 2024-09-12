const URL = "http://localhost:8080";

type StatusResponse = {
  status: string;
};

export async function getStatus(): StatusResponse {
  try {
    const response = await fetch(`${URL}/status`);
    if (!response.ok) {
      throw new Error(`Response status: ${response.status}`);
    }

    return await response.json();
  } catch (error) {
    return {
      status: error.message,
    };
  }
}
