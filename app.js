async function loadBrands() {
  // 根據當前路徑深度決定資料路徑
  let prefix = '';
  if (window.location.pathname.includes('/brands/')) {
    prefix = '../../';
  }
  const res = await fetch(`${prefix}data/brands.json`);
  if (!res.ok) throw new Error('無法載入品牌資料');
  return res.json();
}

function getQuery(name) {
  const params = new URLSearchParams(window.location.search);
  return params.get(name);
}

function getAssetPath(url) {
  if (!url) return '';
  // 如果是絕對路徑或外部連結，直接回傳
  if (url.startsWith('http') || url.startsWith('//') || url.startsWith('data:')) {
    return url;
  }

  let prefix = '';
  if (window.location.pathname.includes('/brands/')) {
    prefix = '../../';
  }

  // 移除可能重複的 ../
  const cleanUrl = url.replace(/^(\.\.\/)+/, '');
  return `${prefix}${cleanUrl}`;
}

function createProductCard(product) {
  const div = document.createElement('div');
  div.className = 'card';
  div.innerHTML = `
    <img src="${getAssetPath(product.image)}" alt="${product.name}" loading="lazy" />
    <div class="card-body">
      <h3>${product.name}</h3>
      <p class="summary">${product.summary}</p>
      <p class="price">$${product.price} · ${product.spec}</p>
    </div>
  `;
  return div;
}

function createBrandCard(brand) {
  const div = document.createElement('div');
  div.className = 'brand-card';
  // 優先使用資料夾結構
  const targetLink = `brands/${brand.slug}/index.html`;

  div.innerHTML = `
    <div class="brand-image">
      <img src="${getAssetPath(brand.logo || brand.heroImage)}" alt="${brand.name}" loading="lazy" />
      <div class="brand-overlay">
        <span class="view-tag">進入品牌</span>
      </div>
    </div>
    <div class="brand-info">
      <h3>${brand.name}</h3>
      <p class="tagline">${brand.tagline}</p>
      <p class="origin">📍 ${brand.origin}</p>
      <div class="brand-badges-mini">
        ${(brand.badges || []).slice(0, 2).map(b => `<span>${b}</span>`).join('')}
      </div>
    </div>
    <a href="${targetLink}" class="brand-link-overlay"></a>
  `;

  return div;
}

function renderHero(brand) {
  const hero = document.querySelector('.lp-hero');
  if (!hero) return;
  if (brand.heroImage) {
    hero.style.setProperty('--hero-bg', `url('${getAssetPath(brand.heroImage)}')`);
  }
  const badgeWrap = document.querySelector('.hero-badges');
  if (badgeWrap) {
    badgeWrap.innerHTML = (brand.badges || []).map(b => `<span class="badge">${b}</span>`).join('');
  }
  const brandTitle = document.querySelector('.hero-brand');
  if (brandTitle) brandTitle.textContent = brand.name || '';
  const tagline = document.querySelector('.hero-tagline');
  if (tagline) tagline.textContent = brand.tagline || '';
  const story = document.querySelector('.hero-story');
  if (story) story.textContent = brand.story || '';
  const ctas = document.querySelectorAll('.hero-cta');
  ctas.forEach(cta => {
    cta.href = brand.cta.link;
    cta.textContent = brand.cta.label || '立即行動';
  });
  const bigText = document.querySelector('.hero-big-text');
  if (bigText) {
    bigText.innerHTML = (brand.heroWords || []).map(w => `<span>${w}</span>`).join('');
  }
}

function renderStory(brand) {
  const text = document.querySelector('.story-text');
  if (text) text.textContent = brand.story;
  const storyImg = document.querySelector('.story-img');
  if (storyImg && brand.heroImage) {
    storyImg.src = getAssetPath(brand.heroImage);
    storyImg.alt = brand.name || '';
  }
}

function renderProducts(brand) {
  const wrap = document.querySelector('.product-grid');
  if (!wrap) return;
  wrap.innerHTML = '';
  (brand.products || []).forEach(p => wrap.appendChild(createProductCard(p)));
}

function renderVideos(brand) {
  const list = document.querySelector('.video-list');
  if (!list || !brand.videos) return;
  list.innerHTML = '';
  brand.videos.forEach(v => {
    const item = document.createElement('div');
    item.className = 'video-item';
    item.innerHTML = `<div>${v.title}</div><a href="${v.url}" target="_blank" rel="noopener">播放影片</a>`;
    list.appendChild(item);
  });
}

function renderFAQ(brand) {
  const faq = document.querySelector('#faq-list');
  if (!faq) return;
  faq.innerHTML = '';
  (brand.faq || []).forEach(item => {
    const div = document.createElement('div');
    div.className = 'faq-item';
    div.innerHTML = `<strong>Q：${item.q}</strong><div>A：${item.a}</div>`;
    faq.appendChild(div);
  });
}

function renderFooterCTA(brand) {
  const footerCta = document.querySelector('#footer-cta');
  if (footerCta) {
    footerCta.href = brand.cta.link;
    footerCta.textContent = brand.cta.label || '立即購買';
    footerCta.style.display = 'block';
  }
}

async function initIndex() {
  const grid = document.getElementById('brand-grid');
  if (!grid) return;

  try {
    const data = await loadBrands();
    grid.innerHTML = '';
    data.brands.forEach(brand => {
      grid.appendChild(createBrandCard(brand));
    });
  } catch (err) {
    console.error(err);
    grid.innerHTML = '<p class="error">載入失敗，請確認 data/brands.json 是否正確生成。</p>';
  }
}

async function initBrand() {
  // 支援 URL query 或路徑推導
  let slug = getQuery('slug');

  // 如果網址路徑包含 brands/xxx/index.html，嘗試推導 slug
  const pathParts = window.location.pathname.split('/');
  const brandsIdx = pathParts.indexOf('brands');
  if (!slug && brandsIdx !== -1 && pathParts[brandsIdx + 1]) {
    slug = pathParts[brandsIdx + 1];
  }

  console.log('Initializing brand page for slug:', slug);
  const data = await loadBrands();
  const brand = data.brands.find(b => b.slug === slug);

  if (!brand) {
    console.warn('Brand not found:', slug);
    const container = document.querySelector('.page-shell');
    if (container) container.innerHTML = '<p style="padding: 100px; text-align: center;">找不到品牌資料。</p>';
    return;
  }

  console.log('Rendering brand:', brand.name);
  document.title = brand.seo?.title || `${brand.name || '品牌頁'}`;

  // Render Logo
  const logoImg = document.getElementById('brand-logo');
  const logoText = document.getElementById('brand-name-text');
  if (logoImg && brand.logo) {
    logoImg.src = getAssetPath(brand.logo);
    logoImg.alt = brand.name || '';
    logoImg.style.display = 'block';
    if (logoText) logoText.style.display = 'none';
  } else if (logoText && brand.name) {
    logoText.textContent = brand.name;
  }

  // Products, Videos, Hero, Story are now handled by SSG for better performance.
  // We only need to render dynamic components like FAQ here.
  renderFAQ(brand);
  renderFooterCTA(brand);
}

window.addEventListener('DOMContentLoaded', () => {
  const page = document.body.dataset.page;
  if (page === 'index') initIndex();
  if (page === 'brand') initBrand();
});
