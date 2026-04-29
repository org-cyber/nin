const express = require('express');
const app = express();
const PORT = process.env.PORT || 8080;

app.use(express.json());
app.use(express.static('static'));

const API_KEY = 'ninbvn_e4ce715a19852fa0b0515f9e';
const PROVIDER_URL = 'https://ninbvnportal.com.ng/api/nin-verification';

app.post('/api/lookup', async (req, res) => {
  const { nin, consent } = req.body;

  if (!nin || !/^\d{11}$/.test(nin)) {
    return res.status(400).json({ error: 'Invalid NIN — exactly 11 digits required' });
  }
  if (consent !== true) {
    return res.status(400).json({ error: 'User consent is required' });
  }

  try {
    const response = await fetch(PROVIDER_URL, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'x-api-key': API_KEY,
      },
      body: JSON.stringify({ nin, consent: true }),
    });

    const data = await response.json();

    if (!response.ok || data.status !== 'success') {
      return res.status(400).json({ error: data.message || 'Verification failed' });
    }

    // ⚠️ The real data is inside data.data (nested)
    const d = data.data?.data;
    if (!d) {
      return res.status(500).json({ error: 'Unexpected response structure' });
    }

    // Assemble address from residence fields
    const addressParts = [
      d.residence_address,
      d.residence_town,
      d.residence_state
    ].filter(Boolean);
    const address = addressParts.length > 0 ? addressParts.join(', ') : '';

    // Return everything the frontend expects
    res.json({
      firstname: d.firstname || '',
      surname: d.surname || '',
      middlename: d.middlename || '',
      fullname: [d.firstname, d.middlename, d.surname].filter(Boolean).join(' '),
      photo: d.photo || '',                          // base64 string, frontend will add data: prefix
      signature: d.signature || '',                  // empty for this provider
      birthdate: d.birthdate || '',
      gender: d.gender || '',
      phone: d.telephoneno || d.phone || '',
      maritalstatus: d.maritalstatus || '',
      birthstate: d.birthstate || '',
      birthlga: d.birthlga || '',
      address: address,
    });
  } catch (err) {
    console.error('Provider error:', err);
    res.status(500).json({ error: 'Something went wrong' });
  }
});

app.listen(PORT, () => {
  console.log(`Server running on http://localhost:${PORT}`);
});

