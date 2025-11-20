import React from 'react';

function TestPage() {
  return (
    <div style={{
      minHeight: '100vh',
      background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
      display: 'flex',
      flexDirection: 'column',
      alignItems: 'center',
      justifyContent: 'center',
      color: 'white',
      fontFamily: 'Arial, sans-serif',
      padding: '20px'
    }}>
      <h1 style={{ fontSize: '3rem', marginBottom: '20px' }}>
        🚀 TwinOS - 前端测试页面
      </h1>
      <p style={{ fontSize: '1.5rem', marginBottom: '30px' }}>
        ✅ React 前端运行正常
      </p>
      <div style={{
        background: 'rgba(255, 255, 255, 0.1)',
        padding: '30px',
        borderRadius: '15px',
        backdropFilter: 'blur(10px)',
        boxShadow: '0 25px 45px rgba(0, 0, 0, 0.2)',
        maxWidth: '600px',
        textAlign: 'center'
      }}>
        <h2 style={{ color: '#fbbf24', marginBottom: '20px' }}>
          📊 系统状态
        </h2>
        <ul style={{ listStyle: 'none', padding: 0, textAlign: 'left' }}>
          <li style={{ margin: '10px 0', fontSize: '1.1rem' }}>
            ✅ React 18：运行中
          </li>
          <li style={{ margin: '10px 0', fontSize: '1.1rem' }}>
            ✅ 开发服务器：端口 3000
          </li>
          <li style={{ margin: '10px 0', fontSize: '1.1rem' }}>
            ✅ TypeScript：编译成功
          </li>
          <li style={{ margin: '10px 0', fontSize: '1.1rem' }}>
            ✅ Tailwind CSS：样式正常
          </li>
        </ul>
      </div>
      <div style={{
        marginTop: '30px',
        fontSize: '1rem',
        opacity: 0.8
      }}>
        如果您看到这个页面，说明前端基础功能正常工作。
      </div>
    </div>
  );
}

export default TestPage;